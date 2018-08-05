package model

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
)

type MongoDBConfig struct {
	Url string
}

type CookbookModel struct {
	dbConf  *MongoDBConfig
	session *mgo.Session
	c       *mgo.Collection
}

type CookbookJson struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Processes   []ProcessJson `json:"processes"`
}

func (cj *CookbookJson) Get() (c lib.Cookbook, err error) {
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

type GeneralProcessJson struct {
	Name    string      `json:"name"`
	Process lib.Process `json:"params"`
}

type ProcessJson bson.Raw

func (pj *ProcessJson) Get() (p lib.Process, err error) {
	gpj := new(GeneralProcessJson)
	err = bson.Unmarshal(pj.Data, gcj)
	switch gpj.Name {
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

	err = bson.Unmarshal(pj.Params.Data, p)
	return
}

func NewCookbookModel(dbConf *MongoDBConfig) *CookbookModel {
	return &CookbookModel{
		dbConf: dbConf,
	}
}

func (m *CookbookModel) ListCookbooks() ([]*lib.Cookbook, error) {
	var csj []CookbookJson
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
	var cj CookbookJson
	if err := m.Connect(); err != nil {
		return nil, err
	}
	if err := m.c.FindId(bson.ObjectIdHex(id)).One(&cj); err != nil {
		return nil, err
	}
	c, err := cj.Get()
	return &c, err
}

func (m *CookbookModel) UpdateCookbook(id string, cookbook *lib.Cookbook) error {
	if err := m.Connect(); err != nil {
		return err
	}
	return m.c.UpdateId(bson.ObjectIdHex(id), cookbook)
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
