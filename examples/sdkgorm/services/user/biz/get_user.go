package biz

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
	"gorm.io/gorm"
)

func (b *biz) GetUserById(ctx context.Context, id int) (*entity.User, error) {
	if user, err := b.userRepo.GetUserById(ctx, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.ErrDB(err)
	} else {
		return user, nil
	}
}
