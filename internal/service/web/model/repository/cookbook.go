package repository

import (
	"GoTuringCoffee/internal/service/lib"
	"GoTuringCoffee/internal/service/web/model/entity"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CookbookRepository cookbook repository
type CookbookRepository struct {
	dbConf     *MongoDBConfig
	db         *mongo.Database
	client     *mongo.Client
	collection *mongo.Collection
}

//NewCookbookRepository Create new cookbook repository
func NewCookbookRepository(dbConf *MongoDBConfig, db *mongo.Database, client *mongo.Client) (model *CookbookRepository, err error) {
	return &CookbookRepository{
		dbConf:     dbConf,
		db:         db,
		client:     client,
		collection: db.Collection(dbConf.Collections.Cookbook),
	}, nil
}

// List List all cookbooks in the database
func (m *CookbookRepository) List(ctx context.Context) ([]lib.Cookbook, error) {
	cookbooksBson := []entity.CookbookBson{}
	err := List(ctx, m.collection, &cookbooksBson)
	if err != nil {
		return nil, err
	}

	cookbooks := make([]lib.Cookbook, 0, len(cookbooksBson))
	for i := range cookbooksBson {
		libCookbook, err := cookbooksBson[i].ToLibModel()
		if err != nil {
			continue
		}

		if libCookbook.ID == lib.DefaultCookbookID {
			continue
		}

		cookbooks = append(cookbooks, libCookbook)
	}

	return cookbooks, nil
}

// CreateDefault Create default cookbook
func (m *CookbookRepository) CreateDefault(ctx context.Context) (lib.Cookbook, error) {
	// convert default cookbook to bson
	defaultCookbook := lib.GenerateDefaultCookbook()
	return m.Update(ctx, defaultCookbook)
}

// Create Create cookbook and save it to mongodb
func (m *CookbookRepository) Create(ctx context.Context, cookbook lib.Cookbook) (lib.Cookbook, error) {
	// Convert lib model to bson model
	cb, err := entity.CreateBsonFromCookbookLibModel(cookbook)
	if err != nil {
		return lib.Cookbook{}, err
	}

	id, err := m.collection.InsertOne(ctx, cb)
	if err != nil {
		return lib.Cookbook{}, err
	}

	idObject := id.InsertedID.(primitive.ObjectID)
	cb.ID = idObject

	newCookbook, err := cb.ToLibModel()
	if err != nil {
		return lib.Cookbook{}, err
	}

	return newCookbook, nil
}

// Update Update cookbook
func (m *CookbookRepository) Update(ctx context.Context, cookbook lib.Cookbook) (lib.Cookbook, error) {
	objectID, err := primitive.ObjectIDFromHex(cookbook.ID)
	if err != nil {
		return cookbook, err
	}

	cb, err := entity.CreateBsonFromCookbookLibModel(cookbook)
	if err != nil {
		return cookbook, err
	}
	cb.UpdatedAt = time.Now().UTC().Unix()

	filter := bson.M{"_id": objectID}
	_, err = m.collection.ReplaceOne(ctx, filter, cb)
	if err != nil {
		return cookbook, err
	}

	return cookbook, nil
}

// GetDefault return default cookbook, if the default cookbook is not existing, create it.
func (m *CookbookRepository) GetDefault(ctx context.Context) (lib.Cookbook, error) {
	cookbook, err := m.Get(ctx, lib.DefaultCookbookID)
	if err == nil {
		return cookbook, nil
	} else {
		return m.CreateDefault(ctx)
	}
}

// Get Get Cookbook
func (m *CookbookRepository) Get(ctx context.Context, id string) (lib.Cookbook, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return lib.Cookbook{}, err
	}

	var cb entity.CookbookBson
	filter := bson.M{"_id": objectID}
	err = m.collection.FindOne(ctx, filter).Decode(&cb)
	if err != nil {
		return lib.Cookbook{}, err
	}

	c, err := cb.ToLibModel()
	return c, err
}

// Delete Delete cookbook
func (m *CookbookRepository) Delete(ctx context.Context, cookbookID string) error {
	objectID, err := primitive.ObjectIDFromHex(cookbookID)
	if err != nil {
		return err
	}

	m.collection.DeleteOne(ctx, bson.M{"_id": objectID})

	return nil
}

// DeleteAll Delete all cookbooks
func (m *CookbookRepository) DeleteAll(ctx context.Context) error {
	_, err := m.collection.DeleteMany(ctx, bson.D{})
	return err
}

// GetAllFieldUnits return all units of process parameters
func (m *CookbookRepository) GetAllFieldUnits() map[string](interface{}) {

	return map[string](interface{}){
		"coordinate": map[string]string{
			"x": "mm",
			"y": "mm",
			"z": "mm",
		},
		"toz": "mm",
		"radius": map[string]string{
			"from": "mm",
			"to":   "mm",
		},
		"cylinder":    "",
		"time":        "s",
		"water":       "ml",
		"temperature": "°C",
	}
}

// GetDefaultProcess return default process
func (m *CookbookRepository) GetAllDefaultProcesses(ctx context.Context) ([]lib.Process, error) {
	defaultCookbook, err := m.GetDefault(ctx)
	if err != nil {
		return []lib.Process{}, err
	}

	return defaultCookbook.Processes, nil
}

// GetProcessNameList return all processes name
func (m *CookbookRepository) GetProcessNameList() []string {
	defaultProcesses := lib.GetDefaultProcesses()
	nameList := make([]string, len(defaultProcesses))

	for index := range defaultProcesses {
		nameList[index] = defaultProcesses[index].Name
		index++
	}
	return nameList
}
