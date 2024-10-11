package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	goservice "github.com/taimaifika/go-sdk"
	"github.com/taimaifika/go-sdk/plugin/simple"
)

func newService() goservice.Service {
	// New service
	service := goservice.New(
		goservice.WithInitRunnable(simple.NewSimplePlugin("simple")),
	)
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

		// Example for Simple Plugin
		type CanGetValue interface {
			GetValue() string
		}
		log.Println(service.MustGet("simple").(CanGetValue).GetValue())

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
