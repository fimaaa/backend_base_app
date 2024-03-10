package application

import "backend_base_app/controller"

type RegistryContract interface {
	controller.Controller
	RunApplication()
}

func Run(rv RegistryContract) {
	if rv != nil {
		rv.RegisterRouter()
		rv.RunApplication()
	}
}
