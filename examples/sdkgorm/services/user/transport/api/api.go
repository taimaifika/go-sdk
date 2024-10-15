package api

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

type Biz interface {
	ListUser(ctx context.Context, limit, offset int) ([]*entity.User, error)
	AddNewUser(ctx context.Context, data *entity.User) error
	GetUserById(ctx context.Context, id int) (*entity.User, error)
	UpdateUserById(ctx context.Context, id int, data *entity.User) error
	DeleteUserById(ctx context.Context, id int) error
}

type api struct {
	biz Biz
}

func NewAPI(biz Biz) *api {
	return &api{biz: biz}
}
