package handler

import (
	"fmt"
	"net/http"

	appContext "league-sim/internal/contexts/appContexts"
	"league-sim/internal/contexts/services"
	"league-sim/internal/models"

	"github.com/labstack/echo/v4"
)

func GetLeagueIds(c echo.Context) error {
	appCtx := c.Request().Context().Value("appContext").(appContext.AppContext)
	appCtx, ok := appCtx.(appContext.AppContext)
	if !ok || appCtx == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	leagueIds, err := appCtx.LeagueRepository().GetLeague()

	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get league IDs")
	}

	return c.JSON(http.StatusOK, leagueIds)
}

func GetLeague(c echo.Context) error {
	leagueId := c.Param("leagueId")

	return c.JSON(http.StatusOK, leagueId)
}

func CreateLeague(c echo.Context) error {
	var body models.CreateLeagueRequest
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	appCtx := c.Request().Context().Value("appContext").(appContext.AppContext)
	appCtx, ok := appCtx.(appContext.AppContext)
	if !ok || appCtx == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	serviceInit := c.Request().Context().Value("services").(services.Service)
	result, err := serviceInit.LeagueService().CreateLeague(body.TeamCount, body.LeagueName)

	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, result)
}

func GetStanding(c echo.Context) error {
	appCtx := c.Request().Context().Value("appContext").(appContext.AppContext)
	appCtx, ok := appCtx.(appContext.AppContext)

	if !ok || appCtx == nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	leagueId := c.Param("leagueId")
	standings, err := appCtx.ActiveLeagueRepository().GetActiveLeaguesStandings(leagueId)

	if err != nil {
		fmt.Println(err)

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get standings")
	}

	return c.JSON(http.StatusOK, standings)
}

func GetFixtures(c echo.Context) error {
	appCtx := c.Request().Context().Value("appContext").(appContext.AppContext)
	appCtx, ok := appCtx.(appContext.AppContext)

	if !ok || appCtx == nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	leagueId := c.Param("leagueId")
	fixtures, err := appCtx.ActiveLeagueRepository().GetActiveLeaguesFixtures(leagueId)

	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get fixtures")
	}

	return c.JSON(http.StatusOK, fixtures)
}

func DeleteLeague(c echo.Context) error {
	leagueId := c.Param("leagueId")

	appCtx := c.Request().Context().Value("appContext").(appContext.AppContext)
	appCtx, ok := appCtx.(appContext.AppContext)
	if !ok || appCtx == nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	err := appCtx.LeagueRepository().DeleteLeague(leagueId)

	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete league")
	}

	return c.JSON(http.StatusOK, leagueId)
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

func GetPredictTable(c echo.Context) error {
	leagueId := c.Param("leagueId")
	service := c.Request().Context().Value("services").(services.Service)
	predictTable, err := service.PredictService().PredictChampionShipSession(leagueId)

	if err != nil {
		fmt.Println("Error predicting championship:", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to predict championship")
	}

	return c.JSON(http.StatusOK, predictTable)
}

func ResetLeague(c echo.Context) error {
	serviceInit := c.Request().Context().Value("services").(services.Service)
	leagueId := c.Param("leagueId")

	err := serviceInit.LeagueService().ResetLeague(leagueId)

	if err != nil {
		fmt.Println("Error resetting league:", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to reset league")
	}

	return c.JSON(http.StatusOK, "League reset successfully")
}
