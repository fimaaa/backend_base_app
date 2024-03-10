package getallmemberv1

import (
	"backend_base_app/domain/entity"
	"context"
)

type Inport interface {
	Execute(ctx context.Context, req entity.BaseReqFind) ([]entity.MemberDataShown, int64, error)
}
