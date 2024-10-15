package biz

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

func (b *biz) ListUser(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return b.userRepo.ListUser(ctx, limit, offset)
}
