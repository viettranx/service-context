package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/component/ginc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GINComponent interface {
	GetPort() int
	GetRouter() *gin.Engine
}

func main() {
	const compId = "gin"

	serviceCtx := sctx.NewServiceContext(
		sctx.WithName("Simple GIN HTTP Service"),
		sctx.WithComponent(ginc.NewGin(compId)),
	)

	if err := serviceCtx.Load(); err != nil {
		log.Fatal(err)
	}

	comp := serviceCtx.MustGet(compId).(GINComponent)

	router := comp.GetRouter()
	router.Use(gin.Recovery(), gin.Logger())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "pong"})
	})

	logger := serviceCtx.Logger("service")

	// Source code from: https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", comp.GetPort()),
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		logger.Infoln("timeout of 5 seconds.")
	}

	_ = serviceCtx.Stop()
	logger.Info("Server exited")
}
