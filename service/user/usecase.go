package user

import (
	"context"
	"sync"

	"github.com/Blackmocca/opentracing-example/models"
)

type UserUsecase interface {
	FetchAll(ctx context.Context, args *sync.Map) ([]*models.User, error)
	FetchAllWithDatabase(ctx context.Context, args *sync.Map) ([]*models.User, error)
}
