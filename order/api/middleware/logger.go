package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

// RequestLogger is a middleware that logs every request
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()

			// Log request information before processing
			fmt.Printf("[REQUEST] %s | %s %s | %s\n",
				time.Now().Format("2006-01-02 15:04:05"),
				req.Method,
				req.URL.Path,
				req.RemoteAddr,
			)

			// Call the next handler
			err := next(c)

			// Calculate request duration
			duration := time.Since(start)

			// Log response information after processing
			fmt.Printf("[RESPONSE] %s | %s %s | %d | %v\n",
				time.Now().Format("2006-01-02 15:04:05"),
				req.Method,
				req.URL.Path,
				res.Status,
				duration,
			)

			return err
		}
	}
}
