package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

var limit = ratelimit.New(1)

func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(c *gin.Context) {
		now := limit.Take()
		log.Print(now.Sub(prev))
		prev = now
	}
}
