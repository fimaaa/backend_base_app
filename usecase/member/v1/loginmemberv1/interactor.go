package loginmemberv1

import (
	"backend_base_app/domain/entity"
	"backend_base_app/shared/dbhelpers"
	"context"
)

type apibaseappmembercreateInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &apibaseappmembercreateInteractor{
		outport: outputPort,
	}
}

func (r *apibaseappmembercreateInteractor) Execute(ctx context.Context, req *entity.MemberAuth) (*entity.MemberDataShown, error) {
	res := &entity.MemberDataShown{}

	err := dbhelpers.WithoutTransaction(ctx, r.outport, func(ctx context.Context) error {

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
