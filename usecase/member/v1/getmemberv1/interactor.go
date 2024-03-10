package getmemberv1

import (
	"backend_base_app/domain/entity"
	"backend_base_app/shared/dbhelpers"
	"context"
)

type apibaseappmembergetallInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &apibaseappmembergetallInteractor{
		outport: outputPort,
	}
}

func (r *apibaseappmembergetallInteractor) Execute(ctx context.Context, id string) (entity.MemberDataShown, error) {
	var response = entity.MemberDataShown{}
	err := dbhelpers.WithoutTransaction(ctx, r.outport, func(ctx context.Context) error {

		res, err := r.outport.FindOneMemberDataById(ctx, id)
		if err != nil {
			return err
		}

		response = *res

		return nil
	})
	return response, err
}
