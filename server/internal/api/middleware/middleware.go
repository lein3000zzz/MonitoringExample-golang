package middleware

import (
	"server/internal/pkg/domain"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func AuthEchoMiddleware(service domain.SessionService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			_, err := service.CheckSession(context.Request().Header)
			if err != nil {
				return context.NoContent(401)
			}

			return next(context)
		}
	}
}

func ErrorLogMiddleware(logger *zap.SugaredLogger) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				logger.Errorf("Handler error for %s %s: %v",
					c.Request().Method,
					c.Request().URL.Path,
					err)
			}
			return err
		}
	}
}
