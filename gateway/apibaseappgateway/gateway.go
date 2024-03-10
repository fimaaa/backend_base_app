package apibaseappgateway

import (
	"backend_base_app/infrastructure/database"
	"fmt"

	cfg "backend_base_app/config/env"
)

type GatewayApiBaseApp struct {
	// *cache.Cache
	database string
	*database.MongoWithTransactionImpl
	*database.MongoWithoutTransactionImpl
	//firebase
	// AuthClientFirebase *auth.Client
	// DbFirebase         *firebaseDb.Ref
	// DBClientFirebase   *firebaseDb.Client
	// DatabaseFirebase   string
}

func NewGateWayApiBaseApp(config cfg.Config) *GatewayApiBaseApp {
	db := database.NewMongoDefault(config)
	dbName := config.GetString("database.mongodb.database")

	fmt.Println("NewGateWayApiBaseApp " + dbName)
	//redis init
	// TODO

	// TODO ADD DEFAULT USER

	return &GatewayApiBaseApp{
		// Cache:                       cacheConnection,
		database:                    dbName,
		MongoWithoutTransactionImpl: database.NewMongoWithoutTransactionImpl(db),
		MongoWithTransactionImpl:    database.NewMongoWithTransactionImpl(db),
		//firebase
		// AuthClientFirebase: authClientFirebase,
		// DbFirebase:       firebaseConnection,
		// DBClientFirebase: dbClientFirebase,
		// DatabaseFirebase: dbFirebase.DatabaseName,
	}
}
