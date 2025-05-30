package handler

import (
	"fmt"
	"league-sim/internal/appContext"
	"league-sim/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetLeagueIds(c echo.Context) error {
	appCtxVal := c.Request().Context().Value("appContext")
	appCtx, ok := appCtxVal.(appContext.AppContext)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	leagueIds, err := appCtx.Adapt.LeagueRepository().GetLeague()
	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get league IDs")
	}

	return c.JSON(http.StatusOK, leagueIds)
}

func CreateLeague(c echo.Context) error {
	var body models.CreateLeagueRequest
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	appCtxVal := c.Request().Context().Value("appContext")
	appCtx, ok := appCtxVal.(appContext.AppContext)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	result, err := appCtx.Service.LeagueService().CreateLeague(body.TeamCount, body.LeagueName)

	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, result)
}

func GetStanding(c echo.Context) error {
	appCtxVal := c.Request().Context().Value("appContext")
	appCtx, ok := appCtxVal.(appContext.AppContext)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	leagueId := c.Param("leagueId")
	standings, err := appCtx.Adapt.ActiveLeagueRepository().GetActiveLeaguesStandings(leagueId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, standings)
}

func GetFixtures(c echo.Context) error {
	appCtxVal := c.Request().Context().Value("appContext")
	appCtx, ok := appCtxVal.(appContext.AppContext)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	leagueId := c.Param("leagueId")
	fixtures, err := appCtx.Adapt.ActiveLeagueRepository().GetActiveLeaguesFixtures(leagueId)

	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get fixtures")
	}

	return c.JSON(http.StatusOK, fixtures)
}

func DeleteLeague(c echo.Context) error {
	leagueId := c.Param("leagueId")

	appCtxVal := c.Request().Context().Value("appContext")
	appCtx, ok := appCtxVal.(appContext.AppContext)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	err := appCtx.Adapt.LeagueRepository().DeleteLeague(leagueId)

	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete league")
	}

	return c.JSON(http.StatusOK, leagueId)
}

func GetPredictTable(c echo.Context) error {
	appCtxVal := c.Request().Context().Value("appContext")
	appCtx, ok := appCtxVal.(appContext.AppContext)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	leagueId := c.Param("leagueId")
	predictTable, err := appCtx.Service.PredictionService().PredictChampionShipSession(leagueId)

	if err != nil {
		fmt.Println("Error predicting championship:", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to predict championship")
	}

	return c.JSON(http.StatusOK, predictTable)
}

func ResetLeague(c echo.Context) error {
	appCtxVal := c.Request().Context().Value("appContext")
	appCtx, ok := appCtxVal.(appContext.AppContext)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	leagueId := c.Param("leagueId")

	err := appCtx.Service.LeagueService().ResetLeague(leagueId)

	if err != nil {
		fmt.Println("Error resetting league:", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to reset league")
	}

	return c.JSON(http.StatusOK, "League reset successfully")
}
