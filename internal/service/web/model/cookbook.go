package model

import (
	"context"
	"fmt"

	"GoTuringCoffee/internal/service/lib"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	Url        string
	Database   string
	Collection string
}

type CookbookModel struct {
	dbConf     *MongoDBConfig
	client     *mongo.Client
	collection *mongo.Collection
}

type CookbookBson struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Processes   []ProcessBson      `bson:"processes"`
}

func (cj *CookbookBson) Get() (c lib.Cookbook, err error) {
	for _, pj := range cj.Processes {
		p, err := pj.Get()
		if err != nil {
			continue
		}
		c.Processes = append(c.Processes, p)
	}
	c.ID = cj.ID.String()
	c.Name = cj.Name
	c.Description = cj.Description
	return
}

type ProcessBson struct {
	Name    string   `bson:"name"`
	Process bson.Raw `bson:"params"`
}

func (pj *ProcessBson) Get() (p lib.Process, err error) {
	switch pj.Name {
	case "Circle":
		p = new(lib.Circle)
	case "Spiral":
		p = new(lib.Spiral)
	case "Polygon":
		p = new(lib.Polygon)
	case "Fixed":
		p = new(lib.Fixed)
	case "Move":
		p = new(lib.Move)
	case "Wait":
		p = new(lib.Wait)
	case "Mix":
		p = new(lib.Mix)
	case "Home":
		p = new(lib.Home)
	default:
		return nil, fmt.Errorf("Not support process '%s'", pj.Name)
	}

	err = bson.Unmarshal(pj.Process, p)
	return
}

func NewCookbookModel(dbConf *MongoDBConfig) *CookbookModel {
	return &CookbookModel{
		dbConf: dbConf,
	}
}

func (m *CookbookModel) ListCookbooks() ([]lib.Cookbook, error) {
	ctx := context.Background()
	if err := m.Connect(); err != nil {
		return nil, err
	}

	fmt.Printf("1\n")

	n, err := m.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	fmt.Printf("2\n")
	if n == 0 {
		return nil, nil
	}

	cookbookBsons := make([]CookbookBson, n)
	cursor, err := m.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	fmt.Printf("3\n")
	for i := 0; cursor.Next(ctx); i++ {
		cursor.Decode(&cookbookBsons[i])
	}

	fmt.Printf("4\n")
	cookbooks := make([]lib.Cookbook, 0, len(cookbookBsons))
	for i := range cookbookBsons {
		var cookbook lib.Cookbook
		if err := bsonToCookbook(&cookbook, &cookbookBsons[i]); err != nil {
			log.Error().Err(err).Msgf("bsonToCookbook")
			continue
		}
		cookbooks = append(cookbooks, cookbook)
	}

	fmt.Printf("5\n")
	fmt.Printf("%v\n", cookbooks)
	return cookbooks, nil
}

func (m *CookbookModel) GetCookbook(id string) (cookbook lib.Cookbook, err error) {

	ctx := context.Background()

	if err := m.Connect(); err != nil {
		return cookbook, err
	}
	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return cookbook, err
	}

	var cookbookBson CookbookBson
	if err := m.collection.FindOne(ctx, bson.M{"_id": bsonID}).Decode(&cookbookBson); err != nil {
		return cookbook, err
	}

	if err := bsonToCookbook(&cookbook, &cookbookBson); err != nil {
		return cookbook, err
	}
	return cookbook, nil
}

func (m *CookbookModel) CreateCookbook(cookbook *lib.Cookbook) error {
	ctx := context.Background()
	if err := m.Connect(); err != nil {
		return err
	}

	var cookbookBson CookbookBson
	if err := cookbookToBson(&cookbookBson, cookbook); err != nil {
		return err
	}
	if _, err := m.collection.InsertOne(ctx, &cookbookBson); err != nil {
		return err
	}
	return nil
}

func (m *CookbookModel) UpdateCookbook(id string, cookbook *lib.Cookbook) error {
	ctx := context.Background()
	if err := m.Connect(); err != nil {
		return err
	}

	var err error
	var bsonID primitive.ObjectID
	var cookbookBson CookbookBson
	if err = cookbookToBson(&cookbookBson, cookbook); err != nil {
		return err
	}

	if bsonID, err = primitive.ObjectIDFromHex(id); err != nil {
		return err
	}

	if _, err = m.collection.UpdateOne(ctx, bson.M{"_id": bsonID}, &cookbookBson); err != nil {
		return err
	}

	return nil
}

func (m *CookbookModel) DeleteCookbook(id string) error {
	if err := m.Connect(); err != nil {
		return err
	}

	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if _, err = m.collection.DeleteOne(context.Background(), bson.M{"_id": bsonID}); err != nil {
		return err
	}
	return nil
}

func (m *CookbookModel) Connect() (err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(m.dbConf.Url))
	if err != nil {
		return err
	}
	if err := client.Connect(context.Background()); err != nil {
		return err
	}

	db := client.Database(m.dbConf.Database)
	collection := db.Collection(m.dbConf.Collection)

	m.client = client
	m.collection = collection
	return
}

func (m *CookbookModel) Disconnect() {
	if m.collection != nil {
		m.client.Disconnect(context.Background())
		m.client = nil
		m.collection = nil
	}
}

func defaultCircle() lib.Process {
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

func defaultSpiral() lib.Process {
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

func defaultPolygon() lib.Process {
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

func defaultFixed() lib.Process {
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

func defaultHome() lib.Process {
	return &lib.Home{}
}

func defaultMove() lib.Process {
	return &lib.Move{
		Coords: lib.Coordinate{
			X: 0,
			Y: 0,
			Z: 200,
		},
	}
}

func defaultWait() lib.Process {
	return &lib.Wait{
		Time: 10,
	}
}

func defaultMix() lib.Process {
	return &lib.Mix{
		Temperature: 80,
	}
}

func getDefaultProcesses() map[string](func() lib.Process) {
	return map[string](func() lib.Process){
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

func getAllDefaultProcesses() map[string](lib.Process) {
	return map[string](lib.Process){
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

func (m *CookbookModel) GetAllFieldUnits() map[string](interface{}) {
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

func (m *CookbookModel) GetDefaultProcess(name string) lib.Process {
	defaultProcesses := getDefaultProcesses()
	if val, ok := defaultProcesses[name]; ok {
		return val()
	} else {
		return nil
	}
}

// GetAllDefaultProcesses return all default processes
func (m *CookbookModel) GetAllDefaultProcesses() map[string](lib.Process) {
	return getAllDefaultProcesses()
}

func (m *CookbookModel) GetProcessNameList() []string {
	defaultProcesses := getDefaultProcesses()
	nameList := make([]string, len(defaultProcesses))

	index := 0
	for key := range defaultProcesses {
		nameList[index] = key
		index++
	}
	return nameList
}

func bsonToCookbook(dst *lib.Cookbook, src *CookbookBson) error {
	dst.ID = src.ID.Hex()
	dst.Name = src.Name
	dst.Description = src.Description
	dst.Processes = make([]lib.Process, len(src.Processes))
	for i, pbson := range src.Processes {
		p := dst.Processes[i]
		switch pbson.Name {
		case "Circle":
			p = new(lib.Circle)
		case "Spiral":
			p = new(lib.Spiral)
		case "Polygon":
			p = new(lib.Polygon)
		case "Fixed":
			p = new(lib.Fixed)
		case "Move":
			p = new(lib.Move)
		case "Wait":
			p = new(lib.Wait)
		case "Mix":
			p = new(lib.Mix)
		case "Home":
			p = new(lib.Home)
		default:
			break
		}

		if err := bson.Unmarshal(pbson.Process, p); err != nil {
			return fmt.Errorf("bson to cookbook: %w", err)
		}
		dst.Processes[i] = p
	}
	return nil
}

func cookbookToBson(dst *CookbookBson, src *lib.Cookbook) error {
	var err error
	if src.ID != "" {
		if dst.ID, err = primitive.ObjectIDFromHex(src.ID); err != nil {
			fmt.Printf(fmt.Sprintf("cookbookToBson: %s %e", src.ID, err))
			return err
		}
	}
	dst.Name = src.Name
	dst.Description = src.Description
	dst.Processes = make([]ProcessBson, len(src.Processes))
	for i, p := range src.Processes {
		pbson := &dst.Processes[i]
		switch p.(type) {
		case *lib.Circle:
			pbson.Name = "Circle"
		case *lib.Spiral:
			pbson.Name = "Spiral"
		case *lib.Polygon:
			pbson.Name = "Polygon"
		case *lib.Fixed:
			pbson.Name = "Fixed"
		case *lib.Move:
			pbson.Name = "Move"
		case *lib.Wait:
			pbson.Name = "Wait"
		case *lib.Mix:
			pbson.Name = "Mix"
		case *lib.Home:
			pbson.Name = "Home"
		default:
			break
		}

		var err error
		pbson.Process, err = bson.Marshal(p)
		if err != nil {
			return err
		}
	}
	return nil
}
