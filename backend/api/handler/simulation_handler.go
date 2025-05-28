package handler

import (
	"fmt"
	"net/http"
	"reflect"

	appContext "league-sim/internal/contexts/appContexts"
	"league-sim/internal/contexts/services"
	"league-sim/internal/models"

	"github.com/labstack/echo/v4"
)

func StartSimulation(c echo.Context) error {
	leagueId := c.Param("leagueId")
	var body models.SimulateLeagueRequest

	if err := c.Bind(&body); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	appCtx := c.Request().Context().Value("appContext").(appContext.AppContext)
	appCtx, ok := appCtx.(appContext.AppContext)

	if !ok || appCtx == nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	service := c.Request().Context().Value("services").(services.Service)

	if service == nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Service context missing")
	}

	result, err := service.SimulationService().Simulation(leagueId, body.PlayAllFixture)

	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start simulation: "+err.Error())
	}

	if reflect.DeepEqual(result, models.SimulationResponse{}) {

		return echo.NewHTTPError(http.StatusNotFound, "No matches to simulate")
	}

	return c.JSON(http.StatusOK, result)
}

func GetMatchResults(c echo.Context) error {
	appCtx := c.Request().Context().Value("appContext").(appContext.AppContext)
	appCtx, ok := appCtx.(appContext.AppContext)

	if !ok || appCtx == nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	leagueId := c.Param("leagueId")
	matchResults, err := appCtx.MatchResultRepository().GetMatchResults(leagueId)

	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get match results: "+err.Error())
	}

	return c.JSON(http.StatusOK, matchResults)
}

func EditMatch(c echo.Context) error {
	var body models.EditMatchResult

	if err := c.Bind(&body); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	leagueId := c.Param("leagueId")
	body.LeagueId = leagueId
	service := c.Request().Context().Value("services").(services.Service)

	err := service.SimulationService().EditMatch(body)
	if err != nil {
		fmt.Println("Error editing match:", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to edit match")
	}

	return c.JSON(http.StatusOK, "Standings updated successfully")
}
