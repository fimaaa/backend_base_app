package creatememberv1

import (
	"backend_base_app/domain/entity"
	"backend_base_app/shared/dbhelpers"
	"backend_base_app/shared/log"
	"backend_base_app/shared/util"
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

func (r *apibaseappmembercreateInteractor) Execute(ctx context.Context, req entity.CreateMemberData) (*entity.MemberDataShown, error) {
	res := &entity.MemberDataShown{}

	err := dbhelpers.WithoutTransaction(ctx, r.outport, func(ctx context.Context) error {

		//automapper
		var memberDataRequest entity.CreateMemberData
		err := util.Automapper(req, &memberDataRequest)
		if err != nil {
			return err
		}
		memberDataObj, err := entity.NewMemberData(memberDataRequest)

		if err != nil {
			return err
		}

		//encrypt password
		password := r.outport.EncryptPassword(ctx, req.Password)
		memberDataObj.Password = password

		log.Info(ctx, util.StructToJson(memberDataObj))

		err = r.outport.CreateMemberData(ctx, *memberDataObj)
		if err != nil {
			return err
		}

		memberShown := memberDataObj.ToShown()
		res = &memberShown

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
