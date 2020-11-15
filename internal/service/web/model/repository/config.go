package repository

import "go.mongodb.org/mongo-driver/mongo"

// SettingsRepository Settings
type SettingsRepository struct {
	Path string
}

//NewSettingsRepository Create new settings repository
func NewSettingsRepository(dbConf *MongoDBConfig, db *mongo.Database, client *mongo.Client) (model *CookbookRepository, err error) {
	return &CookbookRepository{
		dbConf:     dbConf,
		db:         db,
		client:     client,
		collection: db.Collection(dbConf.Collections.Cookbook),
	}, nil
}
