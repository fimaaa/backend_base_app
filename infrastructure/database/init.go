package database

import (
	cfg "backend_base_app/config/env"
	"fmt"

	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSQLiteDefault(config cfg.Config) (db *gorm.DB) {

	db, err := gorm.Open(sqlite.Open("sqlite-database.db"), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	return db
}

func NewPostgresDefault(config cfg.Config) (db *gorm.DB) {
	var err error

	if config.GetString("database.postgresql.host") == "" || config.GetString("database.postgresql.password") == "" || config.GetString("database.postgresql.username") == "" {
		panic(fmt.Errorf("user or password ord databaseName is empty"))
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%v", config.GetString("database.postgresql.host"), config.GetString("database.postgresql.port"), config.GetString("database.postgresql.username"), config.GetString("database.postgresql.database"), config.GetString("database.postgresql.password"), false)

	loggerMode := logger.Silent

	if config.GetBool("database.postgresql.logmode") {
		loggerMode = logger.Info
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(loggerMode),
	})
	if err != nil {
		panic(err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}

	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(10)

	sqlDB.SetConnMaxLifetime(10 * time.Second)

	return db
}

func NewMongoDefault(cfg cfg.Config) *mongo.Client {

	URL := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", cfg.GetString("database.mongodb.username"), cfg.GetString("database.mongodb.password"), cfg.GetString("database.mongodb.host"), cfg.GetString("database.mongodb.port"))

	if cfg.GetString("database.mongodb.username") == "" && cfg.GetString("database.mongodb.password") == "" {
		URL = fmt.Sprintf("mongodb://%s/?authSource=admin", cfg.GetString("database.mongodb.host"))
	}

	if cfg.GetString("database.mongodb.password") != "" {
		password := "xxxxxxxxxxxxxxx"
		maskedPasswordURL := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", cfg.GetString("database.mongodb.username"), password, cfg.GetString("database.mongodb.host"), cfg.GetString("database.mongodb.port"))
		fmt.Printf("\n>>>>>>> MongoDB URI : %s\n\n", maskedPasswordURL)
	} else {
		fmt.Printf("\n>>>>>>> MongoDB URI : %s\n\n", URL)
	}

	setup := options.Client()
	setup.ApplyURI(URL)
	setup.SetDirect(true)
	setup.SetMaxConnecting(1)

	client, err := mongo.NewClient(setup)
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		fmt.Println("connect error >>> ", err)
		panic(err)
	}

	//testing the connections
	existingCollectionNames, err := client.Database(cfg.GetString("database.mongodb.database")).ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}
	fmt.Println("existingCollectionNames >>> ", existingCollectionNames)

	return client
}
