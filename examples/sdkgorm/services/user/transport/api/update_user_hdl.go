package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/services/user/entity"
)

func (api *api) UpdateUserHdl() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		uid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrResBadRequest(err))
			return
		}

		// update user
		var user entity.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrResBadRequest(err))
			return
		}

		if err := api.biz.UpdateUserById(c.Request.Context(), uid, &user); err != nil {
			c.JSON(http.StatusInternalServerError, common.ErrResInternal(err))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(user))
	}
}
