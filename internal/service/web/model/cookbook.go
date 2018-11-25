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
