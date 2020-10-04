package model

import (
	"GoTuringCoffee/internal/service/lib"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbConf.Database)

	cookbookRepo, err := NewCookbookRepository(dbConf, db, client)
	if err != nil {
		return nil, err
	}

	manager = &RepositoryManager{
		dbConf:   dbConf,
		db:       db,
		client:   client,
		Cookbook: cookbookRepo,
	}

	return manager, nil
}

// CookbookRepository cookbook repository
type CookbookRepository struct {
	dbConf     *MongoDBConfig
	db         *mongo.Database
	client     *mongo.Client
	collection *mongo.Collection
}

//NewCookbookRepository Create new cookbook model
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
func (cb *CookbookBson) ToLibModel() (c lib.Cookbook, err error) {
	for _, pj := range cb.Processes {
		p, err := pj.ToLibModel()
		if err != nil {
			continue
		}
		c.Processes = append(c.Processes, *p)
	}
	c.ID = cb.ID.Hex()
	c.Name = cb.Name
	c.Description = cb.Description
	return
}

// ToLibModel Convert process db model(ProcessBson) to lib.process(Process)
func (pb *ProcessBson) ToLibModel() (p *lib.Process, err error) {
	var processImpl lib.ProcessImpl

	switch pb.Name {
	case "Circle":
		processImpl = new(lib.Circle)
	case "Spiral":
		processImpl = new(lib.Spiral)
	case "Polygon":
		processImpl = new(lib.Polygon)
	case "Fixed":
		processImpl = new(lib.Fixed)
	case "Move":
		processImpl = new(lib.Move)
	case "Wait":
		processImpl = new(lib.Wait)
	case "Mix":
		processImpl = new(lib.Mix)
	case "Home":
		processImpl = new(lib.Home)
	default:
		return nil, fmt.Errorf("Not support process '%s'", pb.Name)
	}

	// Convert Bson.Raw to binary
	bson.Unmarshal(pb.Impl, processImpl)
	return &lib.Process{
		ID:        pb.ID,
		Name:      pb.Name,
		CreatedAt: time.Unix(pb.CreatedAt, 0),
		UpdatedAt: time.Unix(pb.UpdatedAt, 0),
		Impl:      processImpl,
	}, nil
}

func (m *CookbookRepository) ConvertToBson(cookbook *lib.Cookbook) (*CookbookBson, error) {
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
			return nil, err
		}

		cb.Processes = append(cb.Processes, pb)
	}

	return &cb, nil
}

// CreateCookbook Create cookbook
func (m *CookbookRepository) Create(ctx context.Context, cookbook *lib.Cookbook) (*lib.Cookbook, error) {
	// Convert lib model to bson model
	cb, err := m.ConvertToBson(cookbook)
	if err != nil {
		return nil, err
	}

	id, err := m.collection.InsertOne(ctx, cb)
	if err != nil {
		return nil, err
	}

	idObject := id.InsertedID.(primitive.ObjectID)
	cb.ID = idObject

	newCookbook, err := cb.ToLibModel()
	if err != nil {
		return nil, err
	}

	return &newCookbook, nil
}

// Update Update cookbook
func (m *CookbookRepository) Update(ctx context.Context, cookbook *lib.Cookbook) error {
	objectID, err := primitive.ObjectIDFromHex(cookbook.ID)
	if err != nil {
		return err
	}

	cb, err := m.ConvertToBson(cookbook)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	m.collection.UpdateOne(ctx, filter, cb)

	return nil
}

// Get Get Cookbook
func (m *CookbookRepository) Get(ctx context.Context, id string) (*lib.Cookbook, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var cb CookbookBson
	filter := bson.M{"_id": objectID}
	err = m.collection.FindOne(ctx, filter).Decode(&cb)
	if err != nil {
		return nil, err
	}

	c, err := cb.ToLibModel()
	return &c, err
}

// Delete Delete cookbook
func (m *CookbookRepository) Delete(ctx context.Context, cookbook *lib.Cookbook) error {
	objectID, err := primitive.ObjectIDFromHex(cookbook.ID)
	if err != nil {
		return err
	}

	m.collection.DeleteOne(ctx, bson.M{"_id": objectID})

	return nil
}

func defaultCircle() lib.ProcessImpl {
	return &lib.Circle{
		Coords: lib.Coordinate{
			X: 0,
			Y: 0,
			Z: 200,
		},
		ToZ: 0,
		Radius: lib.Range{
			From: 0,
			To:   2,
		},
		Cylinder:    5,
		Time:        10,
		Water:       150,
		Temperature: 80,
	}
}

func defaultSpiral() lib.ProcessImpl {
	return &lib.Spiral{
		Coords: lib.Coordinate{
			X: 0,
			Y: 0,
			Z: 200,
		},
		ToZ: 0,
		Radius: lib.Range{
			From: 0,
			To:   2,
		},
		Cylinder:    1,
		Time:        10,
		Water:       150,
		Temperature: 80,
	}
}

func defaultPolygon() lib.ProcessImpl {
	return &lib.Polygon{
		Coords: lib.Coordinate{
			X: 0,
			Y: 0,
			Z: 200,
		},
		ToZ: 0,
		Radius: lib.Range{
			From: 0,
			To:   2,
		},
		Polygon:     2,
		Cylinder:    5,
		Time:        10,
		Water:       150,
		Temperature: 80,
	}
}

func defaultFixed() lib.ProcessImpl {
	return &lib.Fixed{
		Coords: lib.Coordinate{
			X: 0,
			Y: 0,
			Z: 200,
		},
		Time:        10,
		Water:       150,
		Temperature: 80,
	}
}

func defaultHome() lib.ProcessImpl {
	return &lib.Home{}
}

func defaultMove() lib.ProcessImpl {
	return &lib.Move{
		Coords: lib.Coordinate{
			X: 0,
			Y: 0,
			Z: 200,
		},
	}
}

func defaultWait() lib.ProcessImpl {
	return &lib.Wait{
		Time: 10,
	}
}

func defaultMix() lib.ProcessImpl {
	return &lib.Mix{
		Temperature: 80,
	}
}

func getDefaultProcesses() map[string](func() lib.ProcessImpl) {
	return map[string](func() lib.ProcessImpl){
		"Circle":  defaultCircle,
		"Spiral":  defaultSpiral,
		"Polygon": defaultPolygon,
		"Fixed":   defaultFixed,
		"Move":    defaultMove,
		"Wait":    defaultWait,
		"Mix":     defaultMix,
		"Home":    defaultHome,
	}
}

func getAllDefaultProcesses() map[string](lib.ProcessImpl) {
	return map[string](lib.ProcessImpl){
		"Circle":  defaultCircle(),
		"Spiral":  defaultSpiral(),
		"Polygon": defaultPolygon(),
		"Fixed":   defaultFixed(),
		"Move":    defaultMove(),
		"Wait":    defaultWait(),
		"Mix":     defaultMix(),
		"Home":    defaultHome(),
	}
}

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

func (m *CookbookRepository) GetDefaultProcess(name string) lib.ProcessImpl {
	defaultProcesses := getDefaultProcesses()
	if val, ok := defaultProcesses[name]; ok {
		return val()
	} else {
		return nil
	}
}

// GetAllDefaultProcesses return all default processes
func (m *CookbookRepository) GetAllDefaultProcesses() map[string](lib.ProcessImpl) {
	return getAllDefaultProcesses()
}

func (m *CookbookRepository) GetProcessNameList() []string {
	defaultProcesses := getDefaultProcesses()
	nameList := make([]string, len(defaultProcesses))

	index := 0
	for key := range defaultProcesses {
		nameList[index] = key
		index++
	}
	return nameList
}
