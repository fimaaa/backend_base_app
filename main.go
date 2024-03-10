package main

import (
	"backend_base_app/application"
	"backend_base_app/application/registry"
	"flag"
	"fmt"
)

func main() {
	fmt.Printf(" - %s\n", "appName")

	appMap := map[string]func() application.RegistryContract{
		"api_base_app": registry.ApiBaseApp(),
	}

	flag.Parse()

	app, exist := appMap[flag.Arg(0)]
	if exist {
		application.Run(app())
	} else {
		fmt.Println("You may try 'go run main.go <app_name>' :")
		for appName := range appMap {
			fmt.Printf(" - %s\n", appName)
		}
	}
}
