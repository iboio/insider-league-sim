package main

import (
	"league-sim/api"
	"league-sim/config"
	appContext "league-sim/internal/contexts/appContexts"
	"league-sim/internal/contexts/services"
)

func main() {
	config.LoadConfig()

	appCtx, err := appContext.AppContextInit()
	if err != nil {
		panic(err)
	}

	service, err := services.BuildService(appCtx)
	if err != nil {
		panic(err)
	}

	err = api.StartServer(appCtx, service)
	if err != nil {
		panic(err)
	}
	select {}
}
