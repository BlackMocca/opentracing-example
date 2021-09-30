package repository

import (
	"context"
	"sync"
	"time"

	"github.com/Blackmocca/opentracing-example/models"
	"github.com/Blackmocca/opentracing-example/service/user"
	"github.com/opentracing/opentracing-go"
)

type userRepository struct {
}

func NewUserRepository() user.UserRepository {
	return &userRepository{}
}

func (u userRepository) FetchAll(ctx context.Context, args *sync.Map) ([]*models.User, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "FetchAll-Repository")
	defer span.Finish()

	time.Sleep(time.Duration(3 * time.Second))
	users := []*models.User{
		&models.User{
			Id:   "1",
			Name: "ธีรโชค",
		},
		&models.User{
			Id:   "2",
			Name: "มงคล",
		},
	}

	if err := getCover(ctx, users); err != nil {
		return nil, err
	}

	if err := getAddress(ctx, users); err != nil {
		return nil, err
	}

	return users, nil
}
