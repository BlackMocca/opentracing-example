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
	psqlRepo user.UserRepository
}

func NewUserUsecase(userRepo user.UserRepository, psqlRepo user.UserRepository) user.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		psqlRepo: psqlRepo,
	}
}

func (u userUsecase) FetchAll(ctx context.Context, args *sync.Map) ([]*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FetchAll-Usecase")
	defer span.Finish()
	time.Sleep(time.Duration(2 * time.Second))
	return u.userRepo.FetchAll(ctx, args)
}

func (u userUsecase) FetchAllWithDatabase(ctx context.Context, args *sync.Map) ([]*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FetchAllWithdatabase-Usecase")
	defer span.Finish()
	time.Sleep(time.Duration(2 * time.Second))
	return u.psqlRepo.FetchAll(ctx, args)
}
