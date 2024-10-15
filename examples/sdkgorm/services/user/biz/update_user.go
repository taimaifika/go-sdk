package biz

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

// UpdateUserById update user by id
func (b *biz) UpdateUserById(ctx context.Context, id int, user *entity.User) error {
	// update user
	err := b.userRepo.UpdateUserById(ctx, id, user)
	if err != nil {
		return common.ErrDB(err)
	}

	return err
}
