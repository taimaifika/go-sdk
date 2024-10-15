package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	goservice "github.com/taimaifika/go-sdk"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/common"
	"github.com/taimaifika/go-sdk/examples/sdkgorm/composer"
	"github.com/taimaifika/go-sdk/plugin/storage/sdkgorm"
)

func newService() goservice.Service {
	// New service
	service := goservice.New(
		goservice.WithName("sdkgorm"),
		// goservice.WithInitRunnable(sdkgorm.NewGormDB("main.mysql", common.PluginDBMain)),
		// goservice.WithInitRunnable(sdkgocql.NewGcqlDB("gocql", common.PluginDBCassandra)),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main.pg", common.PluginDbPostgres)),
	)
	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start GIN-HTTP service",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize service
		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})

			userApiService := composer.ComposeUserAPIService(service)
			v1 := engine.Group("/v1")
			{
				users := v1.Group("/users")
				{
					// items.POST("", ginitem.CreateItem(service))
					// items.GET("/:id", ginitem.GetItem(service))
					// items.PATCH("/:id", ginitem.PatchItem(service))
					// items.POST("/:id", ginitem.UpdateItem(service))
					// items.DELETE("/:id", ginitem.DeleteItem(service))

					// root => composer => api => services.biz => services.repo
					users.POST("", userApiService.CreateUserHdl())
					users.GET("", userApiService.ListUserHdl())
					users.GET("/:id", userApiService.GetUserHdl())
					users.DELETE("/:id", userApiService.DeleteUserHdl())
					users.POST("/:id", userApiService.UpdateUserHdl())
					users.PATCH("/:id", userApiService.PatchUserHdl())
				}

			}
		})

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
