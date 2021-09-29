package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/Blackmocca/opentracing-example/models"
	"github.com/Blackmocca/opentracing-example/service/user"
	"github.com/opentracing/opentracing-go"
)

type userUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(userRepo user.UserRepository) user.UserRepository {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u userUsecase) FetchAll(ctx context.Context, args *sync.Map) ([]*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FetchAll-Usecase")
	defer span.Finish()
	time.Sleep(time.Duration(2 * time.Second))
	return u.userRepo.FetchAll(ctx, args)
}
