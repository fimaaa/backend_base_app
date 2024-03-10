package apibaseappcontroller

import (
	"backend_base_app/domain/domerror"
	appMiddleware "backend_base_app/lib/wrapper/middleware"
	"backend_base_app/shared/log"
	"backend_base_app/shared/util"
	"backend_base_app/usecase/member/v1/getmemberv1"

	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type userContext struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Rolename string `json:"role"`
}

// authorized is an interceptor
func (r *Controller) authorized(inputPort getmemberv1.Inport) gin.HandlerFunc {

	return func(c *gin.Context) {

		traceID := util.GenerateID()
		ctx := log.Context(c.Request.Context(), traceID)

		authorizationHeader := c.Request.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			newErr := domerror.FailUnmarshalRequestBodyError
			c.AbortWithStatus(http.StatusBadRequest)
			r.Helper.SendBadRequest(c, "Request doesn't contains authorization tokenClaim", newErr, traceID)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		tokenClaim, err := appMiddleware.GetMapClaimByKeyJwtToken(r.Config.GetString("handlingsolution.secret"), tokenString, "jti")
		if err != nil {
			fmt.Println("Error ValidateTokenHandler GetMapClaimByKeyJwtToken jti : ", err)
			if err.Error() == "Token is expired" {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.AbortWithStatus(http.StatusBadRequest)
			}
			r.Helper.SendUnauthorizedError(c, err.Error(), r.Helper.EmptyJsonMap(), traceID)
			return
		}

		c.Set("user", tokenClaim)
		c.Set("tokenstring", tokenString)

		if tokenClaim == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			r.Helper.SendUnauthorizedError(c, "Unauthorized Access", r.Helper.EmptyJsonMap(), traceID)
			return
		}

		tokenClaimStr, ok := tokenClaim.(string)
		if !ok {
			// Type assertion failed, handle the error
			fmt.Println("Failed to convert tokenClaim to string")
			// Handle the error
			return
		}

		fmt.Println("CHECK NORMAL TOKEN ", tokenClaimStr)

		authorized, statusCode, messageResponse := checkAuthorizedAccount(
			ctx, tokenClaimStr, inputPort,
		)

		if !authorized {
			if statusCode == -1 {
				statusCode = http.StatusForbidden
			}
			if messageResponse == "" {
				messageResponse = "Forbidden Access"
			}
			c.AbortWithStatus(statusCode)
			r.Helper.SendUnauthorizedError(c, messageResponse, r.Helper.EmptyJsonMap(), traceID)
			return
		}
		return
	}
}

// authorizedRefreshToken is an interceptor
func (r *Controller) authorizedRefreshToken(inputPort getmemberv1.Inport) gin.HandlerFunc {

	return func(c *gin.Context) {

		traceID := util.GenerateID()
		ctx := log.Context(c.Request.Context(), traceID)

		authorizationHeader := c.Request.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			newErr := domerror.FailUnmarshalRequestBodyError
			c.AbortWithStatus(http.StatusBadRequest)
			r.Helper.SendBadRequest(c, "Request doesn't contains authorization Bearer", newErr, traceID)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		tokenClaim, err := appMiddleware.GetMapClaimByKeyJwtToken(r.Config.GetString("handlingsolution.refresh_token_secret"), tokenString, "jti")
		if err != nil {
			fmt.Println("Error ValidateTokenHandler GetMapClaimByKeyJwtToken jti : ", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			r.Helper.SendUnauthorizedError(c, err.Error(), r.Helper.EmptyJsonMap(), traceID)
			return
		}

		c.Set("user", tokenClaim)
		c.Set("tokenstring", tokenString)

		if tokenClaim == nil {
			c.AbortWithStatus(http.StatusForbidden)
			r.Helper.SendForbiddenError(c, "Forbidden Access", r.Helper.EmptyJsonMap(), traceID)
			return
		}

		tokenClaimStr, ok := tokenClaim.(string)
		if !ok {
			// Type assertion failed, handle the error
			fmt.Println("Failed to convert tokenClaim to string")
			// Handle the error
			return
		}

		fmt.Println("CHECK RERESH_TOKEN ", tokenClaimStr)

		authorized, statusCode, messageResponse := checkAuthorizedAccount(
			ctx, tokenClaimStr, inputPort,
		)

		if !authorized {
			if statusCode == -1 {
				statusCode = http.StatusForbidden
			}
			if messageResponse == "" {
				messageResponse = "Forbidden Access"
			}
			c.AbortWithStatus(statusCode)
			r.Helper.SendUnauthorizedError(c, messageResponse, r.Helper.EmptyJsonMap(), traceID)
			return
		}
		return
	}
}

func checkAuthorizedAccount(
	ctx context.Context,
	tokenClaimStr string,
	inputPort getmemberv1.Inport,
) (bool, int, string) {
	// Unmarshal the tokenClaimStr into a map
	var claimMap map[string]interface{}
	err := json.Unmarshal([]byte(tokenClaimStr), &claimMap)
	if err != nil {
		// Failed to unmarshal, handle the error
		fmt.Println("Failed to unmarshal tokenClaim:", err)
		// Handle the error
		return false, -1, ""
	}

	fmt.Println("TAG ClaimMap >>> ", claimMap)

	deviceId, ok := claimMap["id_device"].(string)
	if !ok {
		fmt.Println("TAG DEVICE ID NOT FOUND")
		// Handle the error if the value is not of the expected type
		// Or assign a default value if needed
		// c.AbortWithStatus(http.StatusInternalServerError)
		// r.Helper.SendNotFoundError(c, "Failed to retrieve authorization id claim", nil, traceID)
		return false, http.StatusInternalServerError, "Failed to retrieve authorization id claim"
	}

	id, ok := claimMap["id"].(string)
	if !ok {
		// Handle the error if the value is not of the expected type
		// Or assign a default value if needed
		// c.AbortWithStatus(http.StatusInternalServerError)
		// r.Helper.SendNotFoundError(c, "Failed to retrieve authorization id claim", nil, traceID)
		return false, http.StatusInternalServerError, "Failed to retrieve authorization id claim"
	}

	courierData, err := inputPort.Execute(ctx, id)

	authorized := true
	statusCode := -1
	messageResponse := ""

	if err == nil {
		fmt.Println("TAG COURIERDATA ", courierData.IsSuspend, " || ", courierData.DeviceId, " == ", deviceId)
		if courierData.IsSuspend == true {
			authorized = false
			statusCode = http.StatusForbidden
			messageResponse = UserSuspended.Error()
		}

		if courierData.DeviceId != deviceId {
			fmt.Println("TAG ClaimMap WHEN COMPARE >>> ", claimMap)
			authorized = false
			statusCode = http.StatusForbidden
			messageResponse = UserLoginInOtherDevice.Error()
		}
	}

	fmt.Println("TAG authorized ", authorized, messageResponse)

	return authorized, statusCode, messageResponse
}

func (r *Controller) superadminAuth(c *gin.Context) error {

	var userCtx userContext

	err := json.Unmarshal([]byte(c.GetString("user")), &userCtx)
	if err != nil {
		return domerror.FailUnmarshalResponseBodyError
	}

	fmt.Printf("User: %s", userCtx.Rolename)
	if userCtx.Rolename != "Superadmin" {
		return InsufficientRole
	}

	return nil
}

const UserSuspended domerror.ErrorType = "ER1006 User is Suspended"
const UserLoginInOtherDevice domerror.ErrorType = "ER1006 Device Already Login in other device"
const InsufficientRole domerror.ErrorType = "ER1009 make sure your role has sufficient authorities"
