package pg

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

func (repo *pgRepo) AddNewUser(ctx context.Context, data *entity.User) error {
	if err := repo.db.Create(data).Error; err != nil {
		return common.ErrCannotCreateEntity("user", err)
	}

	return nil
}
