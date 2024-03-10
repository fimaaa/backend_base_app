package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

//START--------------Transaction-----------------

type MongoWithTransactionImpl struct {
	MongoClient *mongo.Client
}

func NewMongoWithTransactionImpl(c *mongo.Client) *MongoWithTransactionImpl {
	return &MongoWithTransactionImpl{MongoClient: c}
}

func (r *MongoWithTransactionImpl) BeginTransaction(ctx context.Context) (context.Context, error) {

	session, err := r.MongoClient.StartSession()
	if err != nil {
		return nil, err
	}

	sessionCtx := mongo.NewSessionContext(ctx, session)

	err = session.StartTransaction()
	if err != nil {
		panic(err)
	}

	return sessionCtx, nil
}

func (r *MongoWithTransactionImpl) CommitTransaction(ctx context.Context) error {

	err := mongo.SessionFromContext(ctx).CommitTransaction(ctx)
	if err != nil {
		return err
	}

	mongo.SessionFromContext(ctx).EndSession(ctx)

	return nil
}

func (r *MongoWithTransactionImpl) RollbackTransaction(ctx context.Context) error {

	err := mongo.SessionFromContext(ctx).AbortTransaction(ctx)
	if err != nil {
		return err
	}

	mongo.SessionFromContext(ctx).EndSession(ctx)

	return nil
}

//END---------------Transaction-----------------
