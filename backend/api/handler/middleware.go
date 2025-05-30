package handler

import (
	"context"
	"league-sim/internal/appContext"

	"github.com/labstack/echo/v4"
)

func ContextMiddleware(appContext appContext.AppContext) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request().WithContext(context.WithValue(c.Request().Context(), "appContext", appContext))
			c.SetRequest(req)
			return next(c)
		}
	}
}
