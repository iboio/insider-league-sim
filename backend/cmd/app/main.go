package main

import (
	"league-sim/api"
	"league-sim/config"
	"league-sim/internal/appContext"
)

func main() {
	config.LoadConfig()

	appCtx, err := appContext.BuildAppContext()
	if err != nil {
		panic(err)
	}

	err = api.StartServer(appCtx)
	if err != nil {
		panic(err)
	}
	select {}
}
