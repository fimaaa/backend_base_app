package loginmemberv1

import (
	"backend_base_app/domain/entity"
	"context"
)

type Inport interface {
	Execute(ctx context.Context, req *entity.MemberAuth) (*entity.MemberDataShown, error)
}
