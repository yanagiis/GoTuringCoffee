package repository

import (
	"GoTuringCoffee/internal/service/lib"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// GroupRepository group manager
type GroupRepository struct {
	dbConf     *MongoDBConfig
	db         *mongo.Database
	client     *mongo.Client
	collection *mongo.Collection
}

//NewGroupRepository Create new cookbook repository
func NewGroupRepository(dbConf *MongoDBConfig, db *mongo.Database, client *mongo.Client) (model *GroupRepository, err error) {
	return &GroupRepository{
		dbConf:     dbConf,
		db:         db,
		client:     client,
		collection: db.Collection(dbConf.Collections.Cookbook),
	}, nil
}

// List list all groups
func (repo *GroupRepository) List(ctx context.Context) (list []lib.Group, err error) {
	return
}

// GetDefault default group
func (repo *GroupRepository) GetDefault(ctx context.Context) (group lib.Group, err error) {
	return
}

// RemoveCookbook remove cookbook from group
func (repo *GroupRepository) RemoveCookbook(ctx context.Context, group lib.Group, id string) error {
	return nil
}

// RemoveSubGroup remove sub group from group
func (repo *GroupRepository) RemoveSubGroup(ctx context.Context, group lib.Group, subGroupID string) error {
	return nil
}

// Create create Group
func (repo *GroupRepository) Create(ctx context.Context) (group lib.Group, err error) {
	return
}

// Get get group by id
func (repo *GroupRepository) Get(ctx context.Context, id string) (group lib.Group, err error) {
	return
}

// DeleteAll all group except default group
func (repo *GroupRepository) DeleteAll(ctx context.Context) (err error) {
	return
}

// Delete delete group by id
func (repo *GroupRepository) Delete(ctx context.Context, id string) (err error) {
	return
}

// Update update group
func (repo *GroupRepository) Update(ctx context.Context, newGroup lib.Group) (err error) {
	return
}
