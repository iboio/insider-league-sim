package handler

import (
	"context"

	appContext "league-sim/internal/contexts/appContexts"
	"league-sim/internal/contexts/services"

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

func ServiceMiddleware(services services.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, "services", services)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)
			return next(c)
		}
	}
}
