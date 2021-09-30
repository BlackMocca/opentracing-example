package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Blackmocca/opentracing-example/models"
	"github.com/Blackmocca/opentracing-example/service/user"
	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/spf13/cast"
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
		// &models.User{
		// 	Id:   "2",
		// 	Name: "มงคล",
		// },
	}

	if err := u.getCover(ctx, users); err != nil {
		return nil, err
	}

	return users, nil
}

func (u userRepository) getCover(ctx context.Context, users []*models.User) error {
	span := opentracing.SpanFromContext(ctx)
	client := resty.New()
	client.Debug = true
	host := "http://127.0.0.1:3000"

	for index, _ := range users {
		id := users[index].Id
		url := fmt.Sprintf("%s/users/%s/cover", host, id)

		req := client.R().EnableTrace()

		ext.SpanKindRPCClient.Set(span)
		ext.HTTPUrl.Set(span, url)
		ext.HTTPMethod.Set(span, "GET")
		span.Tracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)

		resp, err := req.SetContext(ctx).Get(url)
		if err != nil {
			ext.LogError(span, err)
			panic(err)
		}

		if resp.StatusCode() == http.StatusOK {
			var m = map[string]interface{}{}
			json.Unmarshal(resp.Body(), &m)

			users[index].Cover = cast.ToString(m["cover"])
		}
	}
	return nil
}
