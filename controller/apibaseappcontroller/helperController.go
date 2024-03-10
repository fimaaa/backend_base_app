package apibaseappcontroller

import (
	"backend_base_app/domain/entity"
	"backend_base_app/shared/util"
	"fmt"
)

func (r Controller) CreateMemberToken(
	data entity.MemberDataShown,
) (string, error) {
	//marshal authorizationResult into json

	authorizationResultJson := util.StructToJson(data)

	fmt.Println("CREATE MEMBER TOKEN ", authorizationResultJson)

	//get confidentiality time
	tokenConfidentiality := r.Config.GetInt("api_app_base.token_confidentiality_minute")
	return r.Helper.CreateJwtToken(r.Config.GetString("api_app_base.secret"), string(authorizationResultJson), tokenConfidentiality)
}

func (r Controller) CreateMemberRefreshToken(
	data entity.AuthRefreshToken,
) (string, error) {
	//marshal authorizationResult into json
	refreshTokenJson := util.StructToJson(data)

	fmt.Println("CREATE MEMBER REFRESH TOKEN ", refreshTokenJson)

	refreshTokenConfidentiality := r.Config.GetInt("api_app_base.refresh_token_confidentiality_minute")
	return r.Helper.CreateJwtToken(r.Config.GetString("api_app_base.refresh_token_secret"), string(refreshTokenJson), refreshTokenConfidentiality)
}
