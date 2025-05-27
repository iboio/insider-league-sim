package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"league-sim/api/handler"
	"league-sim/config"
	appContext "league-sim/internal/contexts/appContexts"
	"league-sim/internal/contexts/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer(appCtx appContext.AppContext, services services.Service) error {
	e := echo.New()
	e.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{
					http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions,
				},
				AllowHeaders:     []string{"Content-Type", "Authorization"},
				AllowCredentials: true,
			}))
	buildDir := "public"
	e.Static("/", buildDir)

	e.File("/", filepath.Join(buildDir, "index.html"))
	e.GET(
		"/*", func(c echo.Context) error {
			requestPath := filepath.Join(buildDir, c.Request().URL.Path)

			if _, err := os.Stat(requestPath); err == nil {
				return c.File(requestPath)
			}

			return c.File(filepath.Join(buildDir, "index.html"))
		})

	e.Use(handler.ContextMiddleware(appCtx))
	e.Use(handler.ServiceMiddleware(services))
	v1 := e.Group("/api/v1")

	v1.GET("/league", handler.GetLeagueIds)                           // Get all league IDs
	v1.GET("/league/:leagueId/standing", handler.GetStanding)         // Get league standing by ID
	v1.GET("/league/:leagueId/fixtures", handler.GetFixtures)         // Get fixtures for a league by ID
	v1.GET("/league/:leagueId/predict", handler.GetPredictTable)      // Get simulation results for a league by ID
	v1.GET("/league/:leagueId/matchResults", handler.GetMatchResults) // Get match results for a league by ID

	v1.POST("/league", handler.CreateLeague)                         // Create a new league
	v1.POST("/league/:leagueId/simulation", handler.StartSimulation) // Start a league simulation
	v1.POST("/league/:leagueId/reset", handler.ResetLeague)          // Create fixtures for a league by ID

	v1.DELETE("/league/:leagueId", handler.DeleteLeague) // Delete a league by ID

	v1.PUT("/league/:leagueId", handler.EditMatch) // Update league details by ID

	return e.Start(fmt.Sprintf(":%s", config.HTTPPort))
}
