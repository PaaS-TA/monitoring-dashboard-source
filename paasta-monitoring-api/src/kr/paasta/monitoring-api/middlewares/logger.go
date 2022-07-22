package middlewares

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			logEntry := logrus.NewEntry(logger)
			// request_id를 가져와 logEntry에 셋팅
			id := ctx.Request().Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = ctx.Response().Header().Get(echo.HeaderXRequestID)
			}
			logMap := map[string]interface{}{
				//"request_id":    id,
				"uri":           ctx.Request().RequestURI,
				"method":        ctx.Request().Method,
				"path":          ctx.Request().URL.Path,
				"user_agent":    ctx.Request().UserAgent(),
				"status":        ctx.Response().Status,
			}
			logEntry = logEntry.WithFields(logMap)

			// logEntry를 Context에 저장
			req := ctx.Request()
			ctx.SetRequest(req.WithContext(
				context.WithValue(
					req.Context(),
					"LOG",
					logEntry,
				),
			))
			return next(ctx)
		}
	}
}
