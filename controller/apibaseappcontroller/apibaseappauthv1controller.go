package apibaseappcontroller

import (
	"backend_base_app/domain/domerror"
	"backend_base_app/domain/entity"
	"backend_base_app/shared/log"
	"backend_base_app/shared/util"
	"backend_base_app/usecase/authorization/v1/authmemberv1"
	"backend_base_app/usecase/member/v1/getmemberv1"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ApiBaseAppAuthMember(r *Controller) gin.HandlerFunc {
	var inputPort = authmemberv1.NewUsecase(r.DataSource)

	return func(c *gin.Context) {
		traceID := util.GenerateID()
		ctx := log.Context(c.Request.Context(), traceID)

		var req entity.MemberReqAuth
		if err := c.Bind(&req); err != nil {
			newErr := domerror.FailUnmarshalRequestBodyError
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), newErr, traceID)
			return
		}

		fmt.Println("TAG LOGIN REQUEST => ", req)

		res, err := inputPort.Execute(ctx, req)

		if err != nil {
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), fmt.Sprintf("file err : %s", err.Error()), traceID)
			return
		}

		token, err := r.CreateMemberToken(*res)
		if err != nil {
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), fmt.Sprintf("file err : %s", err.Error()), traceID)
			return
		}
		refreshToken, err := r.CreateMemberRefreshToken(entity.AuthRefreshToken{
			Id:       res.ID,
			DeviceId: res.DeviceId,
		})

		fmt.Println("TAG LOGIN RESPONSE => ", res)

		if err != nil {
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), fmt.Sprintf("file err : %s", err.Error()), traceID)
			return
		}

		finalResponse := entity.MemberResAuth{
			ID:             res.ID,
			Username:       res.Username,
			Fullname:       res.Fullname,
			MemberType:     res.MemberType,
			IsSuspend:      res.IsSuspend,
			CreatedAt:      res.CreatedAt,
			UpdatedAt:      res.UpdatedAt,
			LastLogin:      res.LastLogin,
			TokenBroadcast: res.TokenBroadcast,
			DeviceId:       res.DeviceId,
			PhoneNumber:    res.PhoneNumber,
			Email:          res.Email,
			MemberPhoto:    res.MemberPhoto,
			Token:          token,
			RefreshToken:   refreshToken,
		}

		r.Helper.SendSuccess(c, "Success", finalResponse, traceID)
	}
}

func ApiBaseRefreshAuthMember(r *Controller) gin.HandlerFunc {
	var inputPort = getmemberv1.NewUsecase(r.DataSource)

	return func(c *gin.Context) {
		traceID := util.GenerateID()
		ctx := log.Context(c.Request.Context(), traceID)
		req := entity.AuthRefreshToken{}

		var err error
		//get claim from JWT token
		claim, err := r.Helper.GetJsonClaimFromContext(c)
		if err != nil {
			r.Helper.SendUnauthorizedError(c, err.Error(), err.Error(), traceID)
			return
		}

		err = json.Unmarshal([]byte(claim), &req)
		if err != nil {
			fmt.Sprintf("error unmarshal user token : %s", err.Error())
		}

		res, err := inputPort.Execute(ctx, req.Id)

		if err != nil {
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), fmt.Sprintf("file err : %s", err.Error()), traceID)
			return
		}

		token, err := r.CreateMemberToken(res)
		if err != nil {
			log.Error(ctx, err.Error())
			r.Helper.SendBadRequest(c, err.Error(), fmt.Sprintf("file err : %s", err.Error()), traceID)
			return
		}
		refreshToken, err := r.CreateMemberRefreshToken(entity.AuthRefreshToken{
			Id:       res.ID,
			DeviceId: res.DeviceId,
		})

		finalResponse := entity.MemberResAuth{
			ID:             res.ID,
			Username:       res.Username,
			Fullname:       res.Fullname,
			MemberType:     res.MemberType,
			IsSuspend:      res.IsSuspend,
			CreatedAt:      res.CreatedAt,
			UpdatedAt:      res.UpdatedAt,
			LastLogin:      res.LastLogin,
			TokenBroadcast: res.TokenBroadcast,
			DeviceId:       res.DeviceId,
			PhoneNumber:    res.PhoneNumber,
			Email:          res.Email,
			MemberPhoto:    res.MemberPhoto,
			Token:          token,
			RefreshToken:   refreshToken,
		}

		r.Helper.SendSuccess(c, "Success", finalResponse, traceID)

		return
	}
}
