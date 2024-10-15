package pg

import (
	"context"

	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
	"gorm.io/gorm"
)

func (repo *pgRepo) GetUserById(ctx context.Context, id int) (*entity.User, error) {

	var data entity.User

	if err := repo.db.Where(map[string]interface{}{"id": id}).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &data, nil
}
