package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
)

// Recover returns a middleware that recovers from panics and returns a 500 status code
func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					// Log stack trace
					stack := make([]byte, 4<<10) // 4KB
					length := runtime.Stack(stack, false)

					// You can use your own logger here
					fmt.Printf("[PANIC RECOVER] %v %s\n", r, stack[:length])

					// Return a 500 Internal Server Error response
					err := c.JSON(http.StatusInternalServerError, map[string]interface{}{
						"error": "Internal Server Error",
					})
					if err != nil {
						fmt.Printf("Error sending response: %v\n", err)
					}
				}
			}()

			return next(c)
		}
	}
}
