package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoWithoutTransactionImpl struct {
	MongoClient *mongo.Client
}

func NewMongoWithoutTransactionImpl(c *mongo.Client) *MongoWithoutTransactionImpl {
	return &MongoWithoutTransactionImpl{MongoClient: c}
}

func (r *MongoWithoutTransactionImpl) GetDatabase(ctx context.Context) (context.Context, error) {
	session, err := r.MongoClient.StartSession()
	if err != nil {
		return nil, err
	}

	sessionCtx := mongo.NewSessionContext(ctx, session)

	return sessionCtx, nil
}

func (r *MongoWithoutTransactionImpl) Close(ctx context.Context) error {
	mongo.SessionFromContext(ctx).EndSession(ctx)
	return nil
}
