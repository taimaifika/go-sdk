package biz

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

func (b *biz) AddNewUser(ctx context.Context, data *entity.User) error {
	if err := b.userRepo.AddNewUser(ctx, data); err != nil {
		return err
	}

	return nil
}
