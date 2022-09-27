package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	registerHandler(router)
	server := &http.Server{Addr: ":8080", Handler: router}
	server.RegisterOnShutdown(func() {
		log.Println("shutdown in progress ...")
	})
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		//if err := server.ListenAndServeTLS("testdata/server.pem", "testdata/server.key"); err != nil && err != http.ErrServerClosed {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}

	log.Println("server exiting")
}

func registerHandler(g *gin.Engine) {
	g.Use(leakBucket())
	g.StaticFile("/", "/Users/duanjianing/work/project/godemos/cmd/gin/index.html")
	g.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	g.GET("/stream", func(c *gin.Context) {
		var count int
		c.Stream(func(w io.Writer) bool {
			if count == 10 {
				c.SSEvent("message", "done")
				return false
			}

			time.Sleep(time.Second)
			count++
			c.SSEvent("message", fmt.Sprintf("%d %s", count, time.Now().String()))
			return true
		})
	})

	authorized := g.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "123456",
	}))
	authorized.GET("/version", func(c *gin.Context) { c.String(http.StatusOK, gin.Version) })
}
