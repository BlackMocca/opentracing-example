package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"git.innovasive.co.th/backend/psql"
	"github.com/Blackmocca/opentracing-example/models"
	"github.com/Blackmocca/opentracing-example/orm"
	"github.com/Blackmocca/opentracing-example/service/user"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cast"
)

type psqlUserRepository struct {
	db *psql.Client
}

func NewPsqlUserRepository(client *psql.Client) user.UserRepository {
	return &psqlUserRepository{
		db: client,
	}
}

func (p psqlUserRepository) whereCond(args *sync.Map) []string {
	var conds = []string{}

	if v, ok := args.Load("user_type_id"); ok && v != nil {
		sql := fmt.Sprintf("user_types.id::text = '%s'", cast.ToString(v))
		conds = append(conds, sql)
	}

	return conds
}

func (p psqlUserRepository) FetchAll(ctx context.Context, args *sync.Map) ([]*models.User, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "FetchAll-Repository")
	defer span.Finish()

	time.Sleep(time.Duration(3 * time.Second))

	var conds = p.whereCond(args)
	var where string
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}
	sql := fmt.Sprintf(`
		SELECT 
			%s,
			%s
		FROM users
		JOIN
			user_types
		ON
			users.user_type_id = user_types.id
		%s
	`,
		orm.GetSelector(models.User{}),
		orm.GetSelector(models.UserType{}),
		where,
	)

	rows, err := p.db.GetClient().Queryx(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	joinField := []string{models.FIELD_FK_USER_TYPE}
	users, err := p.orm(ctx, rows, joinField)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (p psqlUserRepository) orm(ctx context.Context, rows *sqlx.Rows, joinField []string) ([]*models.User, error) {
	var users = make([]*models.User, 0)
	var mapper, err = orm.NewRowsScan(rows)
	if err != nil {
		return nil, err
	}

	if mapper.TotalRows() > 0 {
		for _, row := range mapper.RowsValues() {
			var user = new(models.User)
			user, err := orm.OrmUser(user, mapper, row, joinField)
			if err != nil {
				return nil, err
			}
			if user != nil {
				exists, err := orm.IsDuplicateByPK(users, user)
				if err != nil {
					return nil, err
				}
				if !exists {
					users = append(users, user)
				}
			}
		}
	}

	if len(users) > 0 {
		for index, _ := range users {
			if err := orm.OrmUserRelation(users[index], mapper, joinField); err != nil {
				return nil, err
			}
		}

		if err := getCover(ctx, users); err != nil {
			return nil, err
		}

		if err := getAddress(ctx, users); err != nil {
			return nil, err
		}
	}

	return users, nil
}
