package pg

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

func (repo *pgRepo) ListUser(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var data []*entity.User
	if err := repo.db.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
