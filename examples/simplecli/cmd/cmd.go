package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/component/ginc"
	"log"
	"os"
)

func newServiceCtx() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName("simple-gin-http"),
		sctx.WithComponent(ginc.NewGin("gin")),
	)
}

type GINComponent interface {
	GetPort() int
	GetRouter() *gin.Engine
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start GIN-HTTP service",
	Run: func(cmd *cobra.Command, args []string) {
		serviceCtx := newServiceCtx()

		if err := serviceCtx.Load(); err != nil {
			log.Fatal(err)
		}

		comp := serviceCtx.MustGet("gin").(GINComponent)

		router := comp.GetRouter()
		router.Use(gin.Recovery(), gin.Logger())

		if err := router.Run(fmt.Sprintf(":%d", comp.GetPort())); err != nil {
			log.Fatal(err)
		}
	},
}

var outEnvCmd = &cobra.Command{
	Use:   "outenv",
	Short: "Output all environment variables to std",
	Run: func(cmd *cobra.Command, args []string) {
		newServiceCtx().OutEnv()
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
