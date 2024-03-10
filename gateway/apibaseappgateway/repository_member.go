package apibaseappgateway

import (
	"backend_base_app/domain/domerror"
	"backend_base_app/domain/entity"
	"backend_base_app/gateway"
	"backend_base_app/shared/dbhelpers"
	"backend_base_app/shared/log"
	"fmt"
	"time"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateMemberDataRepo interface {
	CreateMemberData(ctx context.Context, obj entity.MemberData) error
	FindOneMemberDataById(ctx context.Context, id string) (*entity.MemberDataShown, error)
	UpdateMemberData(ctx context.Context, memberData entity.MemberDataShown) (*entity.MemberDataShown, error)
	FindAllMemberData(ctx context.Context, req entity.BaseReqFind) ([]*entity.MemberDataShown, int64, error)
	MemberLoginAuthorization(ctx context.Context, obj entity.MemberReqAuth) (*entity.MemberDataShown, error)
}

type memberCollection struct {
	*mongo.Collection
}

func (r GatewayApiBaseApp) getMemberCollection() memberCollection {
	return memberCollection{
		r.MongoWithTransactionImpl.MongoClient.Database(r.database).Collection(entity.CollectionMember),
	}
}

func getFilterKeyword(
	obj entity.MemberDataFind,
	onlySimiliar bool,
) primitive.M {
	//====== execute query using transaction ======
	//count the existing users
	keywordFilter := make([]bson.M, 0)

	if obj.Username != "" {
		keyword := bson.M{"username": obj.Username}
		if onlySimiliar {
			keyword = bson.M{"username": primitive.Regex{Pattern: string(obj.Username), Options: "i"}}
		}
		keywordFilter = append(keywordFilter, keyword)
	}
	if obj.Fullname != "" {
		keyword := bson.M{"fullname": obj.Fullname}
		if onlySimiliar {
			keyword = bson.M{"fullname": primitive.Regex{Pattern: string(obj.Fullname), Options: "i"}}
		}
		keywordFilter = append(keywordFilter, keyword)
	}
	if obj.MemberType != "" {
		keyword := bson.M{"member_type": obj.MemberType}
		keywordFilter = append(keywordFilter, keyword)
	}
	if obj.IsSuspend != nil {
		keyword := bson.M{"is_suspend": obj.IsSuspend}
		keywordFilter = append(keywordFilter, keyword)
	}
	if obj.CreatedAtFrom != nil {
		createdAtTo := time.Now()
		if obj.CreatedAtTo != nil {
			createdAtTo = *obj.CreatedAtTo
		}
		keyword := bson.M{
			"created_at": bson.M{
				"$gte": obj.CreatedAtFrom, // Greater than or equal to date_from
				"$lte": createdAtTo,       // Less than or equal to date_to
			},
		}
		// Append the date filter to the keywordFilter slice
		keywordFilter = append(keywordFilter, keyword)
	}
	if obj.UpdatedAtFrom != nil {
		updatedAtTo := time.Now()
		if obj.UpdatedAtTo != nil {
			updatedAtTo = *obj.UpdatedAtTo
		}
		keyword := bson.M{
			"updated_at": bson.M{
				"$gte": obj.CreatedAtFrom, // Greater than or equal to date_from
				"$lte": updatedAtTo,       // Less than or equal to date_to
			},
		}

		// Append the date filter to the keywordFilter slice
		keywordFilter = append(keywordFilter, keyword)
	}
	if obj.LastLoginFrom != nil {
		lastLoginTo := time.Now()
		if obj.LastLoginTo != nil {
			lastLoginTo = *obj.LastLoginTo
		}
		keyword := bson.M{
			"last_login": bson.M{
				"$gte": obj.LastLoginFrom, // Greater than or equal to date_from
				"$lte": lastLoginTo,       // Less than or equal to date_to
			},
		}

		// Append the date filter to the keywordFilter slice
		keywordFilter = append(keywordFilter, keyword)
	}
	if obj.PhoneNumber != "" {
		keyword := bson.M{"phone_number": obj.PhoneNumber}
		if onlySimiliar {
			keyword = bson.M{"phone_number": primitive.Regex{Pattern: string(obj.PhoneNumber), Options: "i"}}
		}
		keywordFilter = append(keywordFilter, keyword)
	}
	if obj.Email != "" {
		keyword := bson.M{"email": obj.Email}
		if onlySimiliar {
			keyword = bson.M{"email": primitive.Regex{Pattern: string(obj.Email), Options: "i"}}
		}
		keywordFilter = append(keywordFilter, keyword)
	}

	var allCriteria []bson.M
	var criteriaKeyword bson.M
	if len(keywordFilter) > 0 {
		criteriaKeyword = bson.M{"$or": keywordFilter}
		allCriteria = append(allCriteria, criteriaKeyword)
	}

	criteria := bson.M{}
	if len(allCriteria) > 0 {
		criteria = bson.M{"$and": allCriteria}
	}
	return criteria
}

func (coll memberCollection) GetTotalMember(ctx context.Context, obj entity.MemberDataFind, onlySimiliar bool) (int64, error) {
	criteria := getFilterKeyword(obj, onlySimiliar)

	countOpts := options.CountOptions{}
	return coll.CountDocuments(ctx, criteria, &countOpts)
}

func (r GatewayApiBaseApp) CreateMemberData(ctx context.Context, obj entity.MemberData) error {
	log.Info(ctx, "called")

	var err error

	memberCollection := &memberCollection{
		r.MongoWithTransactionImpl.MongoClient.Database(r.database).Collection(entity.CollectionMember),
	}

	count, err := memberCollection.GetTotalMember(
		ctx,
		entity.MemberDataFind{
			Username:    obj.Username,
			PhoneNumber: obj.PhoneNumber,
			Email:       obj.Email,
		},
		false,
	)

	if err != nil {
		log.Error(ctx, err.Error())
		return err
	}
	if count > 0 {
		return error(DataRegistraionHasTaken)
	}

	err = dbhelpers.WithTransaction(ctx, r.MongoWithTransactionImpl, func(dbCtx context.Context) error {
		info, err := memberCollection.InsertOne(ctx, obj)
		log.Info(ctx, "info >>> ", info)
		return err
	})

	return err
}

func (r GatewayApiBaseApp) FindOneMemberDataById(ctx context.Context, id string) (*entity.MemberDataShown, error) {
	log.Info(ctx, "called")

	var (
		resultMemberData entity.MemberDataShown
		err              error
	)

	coll := r.getMemberCollection()
	resCol := coll.FindOne(ctx, bson.M{"id": id})
	err = resCol.Decode(&resultMemberData)
	if err != nil {
		log.Error(ctx, err.Error())

		// error return if data not found
		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("member data not found")
		}

		return nil, err
	}
	fmt.Println("TAG resCol ", resCol)
	fmt.Println("TAG FINDONEBYID ", resultMemberData.CreatedAt)

	return &resultMemberData, nil
}

func (r GatewayApiBaseApp) UpdateMemberData(ctx context.Context, memberData entity.MemberDataShown) (*entity.MemberDataShown, error) {

	log.Info(ctx, "called")

	clientData, err := r.FindOneMemberDataById(ctx, memberData.ID)
	if err != nil || clientData == nil {
		return clientData, err
	}

	memberData.UpdatedAt = time.Now().Local().UTC()

	info, err := r.MongoWithTransactionImpl.SaveOrUpdateByCustomId(ctx, r.database, entity.CollectionMember, memberData.ID, memberData)
	log.Info(ctx, "info >>> ", info)
	if err != nil {
		log.Info(ctx, "error >>> "+err.Error())
	}

	return &memberData, err
}

func (r GatewayApiBaseApp) FindAllMemberData(ctx context.Context, req entity.BaseReqFind) ([]*entity.MemberDataShown, int64, error) {
	log.Info(ctx, "called")

	var (
		err  error
		objs []*entity.MemberDataShown
	)

	coll := r.getMemberCollection()

	findData, _ := req.Value.(entity.MemberDataFind)
	criteria := getFilterKeyword(findData, true)

	findOpts := gateway.BaseReqFindToOptOption(req)

	cursor, err := coll.Find(ctx, criteria, &findOpts)
	if err != nil {
		return nil, 0, err
	}

	if err := cursor.All(ctx, &objs); err != nil {
		return nil, 0, err
	}

	//counting
	count, err := coll.GetTotalMember(ctx, findData, true)

	return objs, count, err
}

func (r GatewayApiBaseApp) MemberLoginAuthorization(ctx context.Context, obj entity.MemberReqAuth) (*entity.MemberDataShown, error) {
	log.Info(ctx, "called")

	var (
		resultMemberDataShown *entity.MemberDataShown
		err                   error
	)

	coll := r.getMemberCollection()

	err = coll.FindOne(ctx, bson.M{"username": obj.Username, "password": obj.Password}).Decode(&resultMemberDataShown)

	if err != nil {
		return resultMemberDataShown, err
	}

	if resultMemberDataShown.IsSuspend {
		err = entity.NewMyError("Account is Suspended")
		return resultMemberDataShown, err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	resultMemberDataShown.LastLogin = time.Now().In(loc)

	fmt.Println("DEVICE ID ", obj.DeviceId)

	if obj.DeviceId != "" {
		resultMemberDataShown.DeviceId = obj.DeviceId
	}
	if obj.TokenBroadcast != "" {
		resultMemberDataShown.TokenBroadcast = obj.TokenBroadcast
	}

	// memberData := entity.MemberData{
	// 	ID:             entity.MemberDataID(resultMemberDataShown.ID),
	// 	Username:       resultMemberDataShown.Username,
	// 	Fullname:       resultMemberDataShown.Fullname,
	// 	MemberType:     resultMemberDataShown.MemberType,
	// 	IsSuspend:      resultMemberDataShown.IsSuspend,
	// 	CreatedAt:      resultMemberDataShown.CreatedAt,
	// 	UpdatedAt:      resultMemberDataShown.UpdatedAt,
	// 	TokenBroadcast: resultMemberDataShown.TokenBroadcast,
	// 	LastLogin:      resultMemberDataShown.LastLogin,
	// 	DeviceId:       resultMemberDataShown.DeviceId,
	// 	PhoneNumber:    resultMemberDataShown.PhoneNumber,
	// 	Email:          resultMemberDataShown.Email,
	// 	MemberPhoto:    resultMemberDataShown.MemberPhoto,
	// }

	//update latest login
	return r.UpdateMemberData(ctx, *resultMemberDataShown)
}

const DataRegistraionHasTaken domerror.ErrorType = "ER1006 data registration has been taken"
