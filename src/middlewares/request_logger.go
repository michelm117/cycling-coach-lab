package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RequestLogger(sugar *zap.SugaredLogger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()
			latency := time.Since(start)

			sugar.Infoln(
				res.Status,
				req.RemoteAddr,
				req.Method,
				req.RequestURI,
				req.Proto,
				req.UserAgent(),
				latency,
				req.Body,
			)

			return nil
		}
	}
}
