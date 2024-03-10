package database

import (
	"backend_base_app/shared/log"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *MongoWithTransactionImpl) PrepareCollection(databaseName string, collectionNames []string) {
	db := r.MongoClient.Database(databaseName)
	ctx := context.Background()

	existingCollectionNames, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		panic(err)
	}

	mapCollName := map[string]int{}
	for _, name := range existingCollectionNames {
		mapCollName[name] = 1
	}

	for _, name := range collectionNames {
		if _, exist := mapCollName[name]; !exist {
			r.createCollection(db.Collection(name), db)
		}
	}
}

func (r *MongoWithTransactionImpl) CreateIndexedTTL(coll *mongo.Collection, dateField string, now time.Time, expireDurationInSecond time.Duration) {
	index := mongo.IndexModel{
		Keys:    bson.M{dateField: 1},
		Options: options.Index().SetExpireAfterSeconds(int32(now.Add(expireDurationInSecond).Unix())),
	}

	_, err := coll.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		panic(err)
	}
}

func (r *MongoWithTransactionImpl) createCollection(coll *mongo.Collection, db *mongo.Database) {
	createCmd := bson.D{{"create", coll.Name()}}
	res := db.RunCommand(context.Background(), createCmd)
	err := res.Err()
	if err != nil {
		panic(err)
	}
}

func (r *MongoWithTransactionImpl) SaveOrUpdate(ctx context.Context, databaseName, collectionName string, id string, data any) (any, error) {

	coll := r.MongoClient.Database(databaseName).Collection(collectionName)

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", data}}
	opts := options.Update().SetUpsert(true)

	result, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%v %v %v", result.UpsertedCount, result.ModifiedCount, result.UpsertedID), nil
}

func (r *MongoWithTransactionImpl) SaveOrUpdateByCustomId(ctx context.Context, databaseName, collectionName string, id string, data any) (any, error) {

	coll := r.MongoClient.Database(databaseName).Collection(collectionName)

	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", data}}
	opts := options.Update().SetUpsert(true)

	result, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%v %v %v", result.UpsertedCount, result.ModifiedCount, result.UpsertedID), nil
}

func (r *MongoWithTransactionImpl) UpdateByCustomId(ctx context.Context, databaseName, collectionName string, id string, data any) (any, error) {

	coll := r.MongoClient.Database(databaseName).Collection(collectionName)

	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", data}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("UpsertedCount: %v , ModifiedCount: %v , UpsertedID: %v", result.UpsertedCount, result.ModifiedCount, result.UpsertedID), nil
}

func (r *MongoWithTransactionImpl) DeleteByCustomId(ctx context.Context, databaseName, collectionName string, id string) (any, error) {

	coll := r.MongoClient.Database(databaseName).Collection(collectionName)

	filter := bson.D{{"id", id}}
	result, err := coll.DeleteOne(ctx, filter)
	log.Info(ctx, "info >>>  ", result)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Deleted >>> %v data", result.DeletedCount), nil
}
