package getallmemberv1

import (
	"backend_base_app/domain/entity"
	"backend_base_app/shared/dbhelpers"
	"backend_base_app/shared/util"
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

func (r *apibaseappmembergetallInteractor) Execute(ctx context.Context, req entity.BaseReqFind) ([]entity.MemberDataShown, int64, error) {
	var response = []entity.MemberDataShown{}
	var totalRecords = int64(-1)
	err := dbhelpers.WithoutTransaction(ctx, r.outport, func(ctx context.Context) error {

		//automapper
		var memberDataRequest entity.BaseReqFind
		err := util.Automapper(req, &memberDataRequest)
		if err != nil {
			return err
		}

		res, count, err := r.outport.FindAllMemberData(ctx, req)
		if err != nil {
			return err
		}

		for _, member := range res {
			response = append(response, *member)
		}

		totalRecords = count

		return nil
	})
	return response, totalRecords, err
}
