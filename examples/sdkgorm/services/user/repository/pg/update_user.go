package pg

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

func (repo *pgRepo) UpdateUserById(ctx context.Context, id int, data *entity.User) error {
	return repo.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Updates(data).Error
}
