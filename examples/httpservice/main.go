package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	goservice "github.com/taimaifika/go-sdk"
	"github.com/taimaifika/go-sdk/sdkcm"
)

func main() {
	service := goservice.New(
		goservice.WithName("demo"),
		goservice.WithVersion("1.0.0"),
	)

	_ = service.Init()

	service.HTTPServer().AddHandler(func(engine *gin.Engine) {
		engine.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"ok":  true,
				"uid": sdkcm.NewUID(1, 1, 1),
			})
		})

		engine.POST("", func(c *gin.Context) {
			type P struct {
				Id sdkcm.UID `json:"id"`
			}

			var p P
			if err := c.ShouldBind(&p); err != nil {
				log.Println(err)
			}

			c.JSON(http.StatusOK, gin.H{
				"local_id": p.Id.GetLocalID(),
				"shard_id": p.Id.GetShardID(),
				"type_id":  p.Id.GetObjectType(),
			})

		})
	})

	_ = service.Start()
}
