package authmemberv1

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

func (r *apibaseappmembercreateInteractor) Execute(ctx context.Context, req entity.MemberReqAuth) (*entity.MemberDataShown, error) {
	response := &entity.MemberDataShown{}

	err := dbhelpers.WithoutTransaction(ctx, r.outport, func(ctx context.Context) error {
		res, err := r.outport.MemberLoginAuthorization(ctx, req)
		if err != nil {
			return err
		}

		response = res

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
