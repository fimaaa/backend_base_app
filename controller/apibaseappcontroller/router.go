package apibaseappcontroller

import (
	cfg "backend_base_app/config/env"
	"backend_base_app/gateway/apibaseappgateway"
	"backend_base_app/shared/helper"
	"backend_base_app/usecase/member/v1/getmemberv1"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Router     gin.IRouter
	Helper     helper.HTTPHelper
	Config     cfg.Config
	DataSource *apibaseappgateway.GatewayApiBaseApp
}

func (r *Controller) handlerAuthMember() gin.HandlerFunc {
	inputPort := getmemberv1.NewUsecase(r.DataSource)
	return r.authorized(inputPort)
}

func (r *Controller) handlerRefreshAuth() gin.HandlerFunc {
	inputPort := getmemberv1.NewUsecase(r.DataSource)
	return r.authorizedRefreshToken(inputPort)
}

func (r *Controller) RegisterRouter() {
	group := r.Router.Group("/api")
	r.RegisterGroupV1(group)
}

func (r *Controller) RegisterGroupV1(groupParent *gin.RouterGroup) {
	group := groupParent.Group("/v1")
	r.RegisterGroupV1Auth(group)
	r.RegisterGroupV1Member(group)
}

func (r *Controller) RegisterGroupV1Auth(groupParent *gin.RouterGroup) {
	group := groupParent.Group("/auth")

	group.POST("/login", ApiBaseAppAuthMember(r))
	group.POST("/refresh", r.handlerRefreshAuth(), ApiBaseRefreshAuthMember(r))
}

func (r *Controller) RegisterGroupV1Member(groupParent *gin.RouterGroup) {
	group := groupParent.Group("/member")

	group.POST("/create", ApiBaseAppMemberCreate(r))
	group.GET("" /* r.handlerAuthMember() ,*/, ApiBaseAppMemberFindAll(r))
	group.GET("/:id", r.handlerAuthMember(), ApiBaseAppMemberFindOne(r))
}
