package apibaseappcontroller

import (
	"backend_base_app/domain/domerror"
	"backend_base_app/domain/entity"
	"backend_base_app/shared/log"
	"backend_base_app/shared/util"
	"backend_base_app/usecase/member/v1/creatememberv1"
	"backend_base_app/usecase/member/v1/getallmemberv1"
	"backend_base_app/usecase/member/v1/getmemberv1"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func ApiBaseAppMemberCreate(r *Controller) gin.HandlerFunc {
	var inputPort = creatememberv1.NewUsecase(r.DataSource)

	return func(c *gin.Context) {
		traceID := util.GenerateID()
		ctx := log.Context(c.Request.Context(), traceID)

		var req entity.CreateMemberData
		if err := c.Bind(&req); err != nil {
			newErr := domerror.FailUnmarshalRequestBodyError
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), newErr, traceID)
			return
		}

		if err := req.ValidateCreate(); err != nil {
			r.Helper.SendBadRequest(c, err.Error(), nil, traceID)
			return
		}

		res, err := inputPort.Execute(ctx, req)

		if err != nil {
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), fmt.Sprintf("file err : %s", err.Error()), traceID)
			return
		}

		r.Helper.SendSuccess(c, "Success", res, traceID)
	}
}

func ApiBaseAppMemberFindAll(r *Controller) gin.HandlerFunc {
	var inputPort = getallmemberv1.NewUsecase(r.DataSource)

	return func(c *gin.Context) {
		traceID := util.GenerateID()
		ctx := log.Context(c.Request.Context(), traceID)

		var req entity.BaseReqFind
		if err := c.BindQuery(&req); err != nil {
			newErr := domerror.FailUnmarshalRequestBodyError
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), newErr, traceID)
			return
		}
		var reqValue entity.MemberDataFind
		c.BindQuery(&reqValue)
		req.Value = reqValue

		sortByParams := make(map[string]interface{})
		// Loop through all query parameters
		for key, value := range c.Request.URL.Query() {
			// Check if the key contains "sort_by_"
			if strings.HasPrefix(key, "sort_by_") {
				trimmedKey := strings.TrimPrefix(key, "sort_by_")
				sortByParams[trimmedKey] = value
			}
		}
		req.SortBy = sortByParams

		res, count, err := inputPort.Execute(ctx, req)

		if err != nil {
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), fmt.Sprintf("file err : %s", err.Error()), traceID)
			return
		}

		finalResponse := req.ToResponse(res, count)

		r.Helper.SendSuccess(c, "Success", finalResponse, traceID)
	}
}

func ApiBaseAppMemberFindOne(r *Controller) gin.HandlerFunc {
	var inputPort = getmemberv1.NewUsecase(r.DataSource)

	return func(c *gin.Context) {
		traceID := util.GenerateID()
		ctx := log.Context(c.Request.Context(), traceID)

		var req string
		if err := c.Bind(&req); err != nil {
			newErr := domerror.FailUnmarshalRequestBodyError
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), newErr, traceID)
			return
		}

		res, err := inputPort.Execute(ctx, req)

		if err != nil {
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), fmt.Sprintf("file err : %s", err.Error()), traceID)
			return
		}

		r.Helper.SendSuccess(c, "Success", res, traceID)
	}
}
