package authmemberv1

import (
	"backend_base_app/domain/service"
	"backend_base_app/gateway/apibaseappgateway"
	"backend_base_app/shared/dbhelpers"
)

type Outport interface {
	service.EncryptPasswordService
	apibaseappgateway.CreateMemberDataRepo
	dbhelpers.WithoutTransactionDB
}
