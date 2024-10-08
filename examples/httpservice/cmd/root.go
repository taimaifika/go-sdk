package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	goservice "github.com/taimaifika/go-sdk"
)

func newService() goservice.Service {
	// New service
	service := goservice.New(
		goservice.WithName("go-sdk"),
	)

	fmt.Println("Service Name:", service.Name())

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start whoami service",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize service
		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			// Health check
			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
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
