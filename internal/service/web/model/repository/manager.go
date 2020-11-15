package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConfig mongodb configuration
type MongoDBConfig struct {
	URL         string
	Database    string
	Collections struct {
		Cookbook string
	}
}

// RepositoryManager Repository manager
type RepositoryManager struct {
	dbConf *MongoDBConfig
	db     *mongo.Database
	client *mongo.Client

	Cookbook *CookbookRepository
}

// NewRepositoryManager create repository manager
func NewRepositoryManager(parentContext context.Context, dbConf *MongoDBConfig) (manager *RepositoryManager, err error) {
	ctx, cancel := context.WithTimeout(parentContext, 10*time.Second)
	options := options.Client()
	options.SetMaxPoolSize(64)
	options.ApplyURI(dbConf.URL)
	defer cancel()

	client, err := mongo.Connect(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect() failed: %v", err.Error())
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Ping mongodb failed: %v", err.Error())
	}

	db := client.Database(dbConf.Database)

	cookbookRepo, err := NewCookbookRepository(dbConf, db, client)
	if err != nil {
		return nil, fmt.Errorf("Can't create cookbook repository: %v", err.Error())
	}

	manager = &RepositoryManager{
		dbConf:   dbConf,
		db:       db,
		client:   client,
		Cookbook: cookbookRepo,
	}

	return manager, nil
}
