package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

func (api *api) CreateUserHdl() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var data entity.User

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrResBadRequest(err))
			return
		}

		if err := api.biz.AddNewUser(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrResInternal(err))
			return
		}

		c.JSON(http.StatusCreated, common.SimpleSuccessResponse(data))
	}
}
