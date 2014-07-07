package domain

import (
	"github.com/gin-gonic/gin"
)

func Bind(host string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Req.Host != host {
			c.Abort(404)
		}
	}
}

func Logger() HandlerFunc {
	channel := make(chan logInfo, 32)
	go func() {
		logger := log.New(os.Stdout, "", 0)
		for info := range channel {
			var color string
			code := info.status
			switch {
			case code >= 200 && code <= 299:
				color = green
			case code >= 300 && code <= 399:
				color = white
			case code >= 400 && code <= 499:
				color = yellow
			default:
				color = red
			}
			//			logger.Printf("[GIN] %v |%s %3d %s| %12v | %3.1f%% | %3s | %s\n",

			logger.Printf("[GIN] %v |%s %3d %s| %12v | %3s | %s\n",
				info.time.Format("2006/01/02 - 15:04:05"),
				color, info.status, reset,
				info.latency,
				info.req.Method, info.req.URL.Path,
			)

			// Calculate resolution time
			if len(info.errors) > 0 {
				fmt.Println(info.errors.String())
			}
		}
	}()
	return func(c *Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		latency := time.Since(start)
		channel <- logInfo{
			time:    start,
			status:  c.Writer.Status(),
			latency: latency,
			req:     c.Req,
			errors:  c.Errors,
		}
	}
}
