package biz

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

type UserRepository interface {
	AddNewUser(ctx context.Context, data *entity.User) error
	ListUser(ctx context.Context, limit, offset int) ([]*entity.User, error)
	GetUserById(ctx context.Context, id int) (*entity.User, error)
	DeleteUserById(ctx context.Context, id int) error
	UpdateUserById(ctx context.Context, id int, data *entity.User) error
}

type biz struct {
	userRepo UserRepository
}

func NewBiz(userRepo UserRepository) *biz {
	return &biz{userRepo: userRepo}
}
