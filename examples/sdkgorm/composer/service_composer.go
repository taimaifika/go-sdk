package composer

import (
	"github.com/gin-gonic/gin"
	goservice "github.com/taimaifika/go-sdk"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/transport/api"
	"gorm.io/gorm"

	userBiz "github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/biz"
	userSqlRepo "github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/repository/pg"
)

type UserService interface {
	CreateUserHdl() func(ctx *gin.Context)
	ListUserHdl() func(ctx *gin.Context)
	GetUserHdl() func(ctx *gin.Context)
	DeleteUserHdl() func(ctx *gin.Context)
	UpdateUserHdl() func(ctx *gin.Context)
	PatchUserHdl() func(ctx *gin.Context)
}

// ComposeUserAPIService compose user api service
func ComposeUserAPIService(serviceCtx goservice.ServiceContext) UserService {
	db := serviceCtx.MustGet(common.PluginDbPostgres).(*gorm.DB)

	userRepo := userSqlRepo.NewPostgresRepository(db)
	biz := userBiz.NewBiz(userRepo)
	userService := api.NewAPI(biz)

	return userService
}
