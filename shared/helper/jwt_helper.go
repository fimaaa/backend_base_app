package helper

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JwtClaims struct {
	// add as necessary
	jwt.StandardClaims
}

// CreateJwtToken ...
func (u HTTPHelper) CreateJwtToken(secret string, apiToken string, confidentialMinute int) (string, error) {
	claims := JwtClaims{
		jwt.StandardClaims{
			Id:        apiToken,
			ExpiresAt: time.Now().Add(time.Duration(confidentialMinute) * time.Minute).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, err
}

func (u HTTPHelper) GetJwtClaims(c *gin.Context) jwt.MapClaims {
	user, _ := c.Get("user")
	if user == nil {
		return nil
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	return claims
}

func (u HTTPHelper) GetJwtClaim(c *gin.Context, key string) interface{} {
	claims := u.GetJwtClaims(c)

	return claims[key]
}

func (u HTTPHelper) GetMemberCode(c *gin.Context, coreUrl string) (string, error) {
	var (
	//err error
	)

	apiToken := u.GetJwtClaim(c, "jti")
	if apiToken == nil {
		return "", errors.New("There is no API Token in JWT claims.")
	}

	//url := str.Replacer(coreUrl, strings.NewReplacer("{apiToken}", apiToken.(string)))
	//// get member code by api token
	//req := client.Info{
	//	Url:    url,
	//	Method: "GET",
	//}

	type result struct {
		Code     int                    `json:"code"`
		CodeType string                 `json:"code_type"`
		Message  string                 `json:"message"`
		Data     map[string]interface{} `json:"data"`
	}

	resp := result{}
	//if err = req.Dispatch(&resp); err != nil {
	//	return "", err
	//}
	//if len(resp.Data) == 0 {
	//	return "", errors.New("Member not found")
	//}

	return resp.Data["member_code"].(string), nil
}

func (u HTTPHelper) GetMemberAPIToken(c *gin.Context) (string, error) {
	apiToken := u.GetJwtClaim(c, "jti")
	if apiToken == nil {
		return "", errors.New("There is no API Token in JWT claims.")
	}

	return apiToken.(string), nil
}

func (u HTTPHelper) GetJsonClaimFromToken(c *gin.Context) (string, error) {
	claim := u.GetJwtClaim(c, "jti")
	if claim == nil {
		return "", errors.New("There is no API Token in JWT claims.")
	}

	return claim.(string), nil
}

func (u HTTPHelper) GetJsonClaimFromContext(c *gin.Context) (string, error) {
	userFromContext, _ := c.Get("user")
	if userFromContext == nil {
		return "", fmt.Errorf("token string not exist in user context")
	}

	return userFromContext.(string), nil
}
