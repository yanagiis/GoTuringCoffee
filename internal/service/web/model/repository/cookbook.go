package repository

import (
	"GoTuringCoffee/internal/service/lib"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	query := bson.M{}
	findOptions := options.Find()
	findOptions.Sort = bson.M{
		"created_at": -1,
	}

	cookbooksResult, err := m.collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}

	cookbooksBson := []CookbookBson{}
	err = cookbooksResult.All(ctx, &cookbooksBson)
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

// ProcessBson Process data structure
type ProcessBson struct {
	ID        string   `bson:"id, omitempty"`
	Name      string   `bson:"name"`
	Impl      bson.Raw `bson:"process_impl"`
	CreatedAt int64    `bson:"created_at"`
	UpdatedAt int64    `bson:"updated_at"`
}

// CookbookBson Cookbook data structure
type CookbookBson struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Tags        []string           `bson:"tags,omitempty"`
	Notes       []string           `bson:"notes,omitempty"`
	Processes   []ProcessBson      `bson:"processes,omitempty"`
	CreatedAt   int64              `bson:"created_at"`
	UpdatedAt   int64              `bson:"updated_at"`
}

// SetID Convert the hex string to ObjectID
func (cb CookbookBson) SetID(id string) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		return
	}

	cb.ID = oid
}

// ToLibModel Convert db model(CookbookBson) to lib.Cookbook(Cookbook)
func (cb CookbookBson) ToLibModel() (c lib.Cookbook, err error) {
	for _, pj := range cb.Processes {
		p, err := pj.ToLibModel()
		if err != nil {
			continue
		}
		c.Processes = append(c.Processes, p)
	}
	c.ID = cb.ID.Hex()
	c.Name = cb.Name
	c.Description = cb.Description
	return
}

// ToLibModel Convert process db model(ProcessBson) to lib.process(Process)
func (pb *ProcessBson) ToLibModel() (p lib.Process, err error) {
	processImpl, err := lib.NewProcessImpl(pb.Name)
	if err != nil {
		return lib.Process{}, err
	}

	// Convert Bson.Raw to binary
	bson.Unmarshal(pb.Impl, processImpl)
	return lib.Process{
		ID:        pb.ID,
		Name:      pb.Name,
		CreatedAt: time.Unix(pb.CreatedAt, 0),
		UpdatedAt: time.Unix(pb.UpdatedAt, 0),
		Impl:      processImpl,
	}, nil
}

// ConvertToBson Convert cookbook model to bson for mongodb
func (m *CookbookRepository) ConvertToBson(cookbook lib.Cookbook) (CookbookBson, error) {
	var cb CookbookBson

	cb.Name = cookbook.Name
	cb.Description = cookbook.Description
	cb.CreatedAt = cookbook.CreatedAt.Unix()
	cb.UpdatedAt = cookbook.UpdatedAt.Unix()
	cb.Tags = cookbook.Tags
	cb.Notes = cookbook.Notes

	for _, p := range cookbook.Processes {
		var pb ProcessBson
		var err error

		pb.ID = p.ID
		pb.Name = p.Name
		pb.CreatedAt = p.CreatedAt.Unix()
		pb.UpdatedAt = p.UpdatedAt.Unix()

		pb.Impl, err = bson.Marshal(p.Impl)
		if err != nil {
			return CookbookBson{}, err
		}

		cb.Processes = append(cb.Processes, pb)
	}

	return cb, nil
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
	cb, err := m.ConvertToBson(cookbook)
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

	cb, err := m.ConvertToBson(cookbook)
	if err != nil {
		return cookbook, err
	}
	cb.UpdatedAt = time.Now().UTC().Unix()

	filter := bson.M{"_id": objectID}
	m.collection.UpdateOne(ctx, filter, cb)

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

	var cb CookbookBson
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
