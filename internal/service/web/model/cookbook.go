package model

import (
	"fmt"

	"GoTuringCoffee/internal/service/lib"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MongoDBConfig struct {
	Url string
}

type CookbookModel struct {
	dbConf  *MongoDBConfig
	session *mgo.Session
	c       *mgo.Collection
}

type CookbookBson struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Processes   []ProcessBson `bson:"processes"`
}

func (cj *CookbookBson) Get() (c lib.Cookbook, err error) {
	for _, pj := range cj.Processes {
		p, err := pj.Get()
		if err != nil {
			continue
		}
		c.Processes = append(c.Processes, p)
	}
	c.ID = cj.ID
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

	err = bson.Unmarshal(pj.Process.Data, p)
	return
}

func NewCookbookModel(dbConf *MongoDBConfig) *CookbookModel {
	return &CookbookModel{
		dbConf: dbConf,
	}
}

func (m *CookbookModel) ListCookbooks() ([]*lib.Cookbook, error) {
	var csj []CookbookBson
	var cs []*lib.Cookbook
	if err := m.Connect(); err != nil {
		return cs, err
	}
	if err := m.c.Find(nil).All(&csj); err != nil {
		return nil, err
	}
	fmt.Printf("Get %d cookbooks from db", len(csj))
	for _, cj := range csj {
		c, _ := cj.Get()
		cs = append(cs, &c)
	}
	return cs, nil
}

func (m *CookbookModel) GetCookbook(id string) (*lib.Cookbook, error) {
	var cj CookbookBson
	if err := m.Connect(); err != nil {
		return nil, err
	}
	if err := m.c.FindId(bson.ObjectIdHex(id)).One(&cj); err != nil {
		return nil, err
	}
	c, err := cj.Get()
	return &c, err
}

func (m *CookbookModel) CreateCookbook(cookbook *lib.Cookbook) error {
	var cb CookbookBson
	if err := m.Connect(); err != nil {
		return err
	}

	cb.Name = cookbook.Name
	cb.Description = cookbook.Description
	for _, p := range cookbook.Processes {
		var pb ProcessBson
		var err error
		switch p.(type) {
		case *lib.Circle:
			pb.Name = "Circle"
		case *lib.Spiral:
			pb.Name = "Spiral"
		case *lib.Fixed:
			pb.Name = "Fixed"
		case *lib.Move:
			pb.Name = "Move"
		case *lib.Wait:
			pb.Name = "Wait"
		case *lib.Mix:
			pb.Name = "Mix"
		case *lib.Home:
			pb.Name = "Home"
		}
		pb.Process.Data, err = bson.Marshal(p)
		if err != nil {
			return err
		}

		cb.Processes = append(cb.Processes, pb)
	}

	return m.c.Insert(cb)
}

func (m *CookbookModel) UpdateCookbook(id string, cookbook *lib.Cookbook) error {
	var cb CookbookBson
	if err := m.Connect(); err != nil {
		return err
	}

	cb.ID = cookbook.ID
	cb.Name = cookbook.Name
	cb.Description = cookbook.Description
	for _, p := range cookbook.Processes {
		var pb ProcessBson
		var err error
		switch p.(type) {
		case *lib.Circle:
			pb.Name = "Circle"
		case *lib.Spiral:
			pb.Name = "Spiral"
		case *lib.Fixed:
			pb.Name = "Fixed"
		case *lib.Move:
			pb.Name = "Move"
		case *lib.Wait:
			pb.Name = "Wait"
		case *lib.Mix:
			pb.Name = "Mix"
		case *lib.Home:
			pb.Name = "Home"
		}
		pb.Process.Data, err = bson.Marshal(p)
		if err != nil {
			return err
		}

		cb.Processes = append(cb.Processes, pb)
	}

	return m.c.UpdateId(bson.ObjectIdHex(id), cb)
}

func (m *CookbookModel) DeleteCookbook(id string) error {
	if err := m.Connect(); err != nil {
		return err
	}
	return m.c.RemoveId(bson.ObjectIdHex(id))
}

func (m *CookbookModel) Connect() (err error) {
	if m.session == nil {
		if m.session, err = mgo.Dial(m.dbConf.Url); err != nil {
			return
		}
	}
	if m.c == nil {
		m.c = m.session.DB("turing-coffee").C("cookbooknew")
	}
	return
}

func (m *CookbookModel) Disconnect() {
	if m.session != nil {
		m.session.Close()
		m.session = nil
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
