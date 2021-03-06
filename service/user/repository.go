package user

import (
	"context"
	"sync"

	"github.com/Blackmocca/opentracing-example/models"
)

type UserRepository interface {
	FetchAll(ctx context.Context, args *sync.Map) ([]*models.User, error)
}
