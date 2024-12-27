package xhttp

import (
	"net/http"
	"strings"
	"time"

	"github.com/aronfan/plat.mini/xlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func MiddlewareMonitor(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path

	c.Next()

	end := time.Now()
	latency := end.Sub(start)
	status := c.Writer.Status()
	xlog.Debug("Request information",
		zap.String("path", path),
		zap.Int("status", status),
		zap.Duration("latency", latency),
	)
}

func MiddlewareCORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	} else {
		c.Next()
	}
}

func MiddlewareUnity(c *gin.Context) {
	path := c.Request.URL.Path

	if strings.HasSuffix(path, ".br") {
		c.Header("Content-Encoding", "br")
		c.Header("Content-Type", "application/wasm")
	}

	if strings.HasSuffix(path, ".js.gz") {
		c.Header("Content-Encoding", "gzip")
		c.Header("Content-Type", "application/javascript")
	}

	if strings.HasSuffix(path, ".wasm.gz") {
		c.Header("Content-Encoding", "gzip")
		c.Header("Content-Type", "application/wasm")
	}

	if strings.HasSuffix(path, ".data.gz") {
		c.Header("Content-Encoding", "gzip")
		c.Header("Content-Type", "application/gzip")
	}

	if strings.HasSuffix(path, ".symbols.json.gz") {
		c.Header("Content-Encoding", "gzip")
		c.Header("Content-Type", "application/json")
	}

	c.Next()
}
