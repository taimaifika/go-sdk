package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
)

func (api *api) DeleteUserHdl() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		uid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrResBadRequest(err))
			return
		}

		err = api.biz.DeleteUserById(c.Request.Context(), uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrResInternal(err))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(nil))
	}
}
