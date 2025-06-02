package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"league-sim/internal/appContext"
	"league-sim/internal/models"
	"net/http"
	"reflect"
)

func StartSimulation(c echo.Context) error {
	leagueId := c.Param("leagueId")
	var body models.SimulateLeagueRequest

	if err := c.Bind(&body); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	appCtxVal := c.Request().Context().Value("appContext")
	appCtx, ok := appCtxVal.(appContext.AppContext)

	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	result, err := appCtx.Service.SimulationService().Simulation(leagueId, body.PlayAllFixture)

	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start simulation: "+err.Error())
	}

	if reflect.DeepEqual(result, models.SimulationResponse{}) {

		return echo.NewHTTPError(http.StatusNotFound, "No matches to simulate")
	}

	return c.JSON(http.StatusOK, result)
}

func EditMatch(c echo.Context) error {
	var body models.EditMatchResult

	if err := c.Bind(&body); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	appCtxVal := c.Request().Context().Value("appContext")
	appCtx, ok := appCtxVal.(appContext.AppContext)

	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	leagueId := c.Param("leagueId")
	body.LeagueId = leagueId

	err := appCtx.Service.SimulationService().EditMatch(body)

	if err != nil {
		fmt.Println("Error editing match:", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to edit match")
	}

	return c.JSON(http.StatusOK, "Standings updated successfully")
}
