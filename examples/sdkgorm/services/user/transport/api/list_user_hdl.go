package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
)

func (api *api) ListUserHdl() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		users, err := api.biz.ListUser(c.Request.Context(), 0, 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrResInternal(err))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(users))
	}
}
