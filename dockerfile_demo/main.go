package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		hostName, _ := os.Hostname()
		c.JSON(200, gin.H{
			"message": "pong",
			"host":    hostName,
		})
	})
	r.GET("/heartbeat", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version": "2.0.0",
		})
	})

	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
