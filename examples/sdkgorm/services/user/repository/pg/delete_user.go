package pg

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

func (repo *pgRepo) DeleteUserById(ctx context.Context, id int) error {
	if err := repo.db.Delete(&entity.User{}, id).Error; err != nil {
		return common.ErrCannotDeleteEntity("user", err)
	}

	return nil
}
