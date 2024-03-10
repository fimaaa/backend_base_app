package registry

import (
	"backend_base_app/application"
	cfg "backend_base_app/config/env"
	"backend_base_app/controller"
	"backend_base_app/controller/apibaseappcontroller"
	"backend_base_app/gateway/apibaseappgateway"
	"backend_base_app/infrastructure/server"
)

type baseapp struct {
	*server.GinHTTPHandler
	controller.Controller
}

func ApiBaseApp() func() application.RegistryContract {
	return func() application.RegistryContract {
		//register config
		config := cfg.NewViperConfig()

		// register webhook
		// TODO future

		httpHandler := server.NewGinHTTPHandlerDefault(config.GetString("api_app_base.api_url"))

		// datasource := apibaseappgateway.NewGateWayApiBaseApp(config)

		return &baseapp{
			GinHTTPHandler: &httpHandler,
			Controller: &apibaseappcontroller.Controller{
				Router:     httpHandler.Router,
				Config:     config,
				DataSource: apibaseappgateway.NewGateWayApiBaseApp(config),
			},
		}
	}
}
