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

type Cookbook struct {
	dbConf  *MongoDBConfig
	session *mgo.Session
	c       *mgo.Collection
}

func NewCookbook(dbConf *MongoDBConfig) *Cookbook {
	return &Cookbook{
		dbConf: dbConf,
	}
}

func (m *Cookbook) ListCookbooks() ([]lib.Cookbook, error) {
	var cookbooks []lib.Cookbook
	if err := m.Connect(); err != nil {
		return nil, err
	}
	if err := m.c.Find(nil).All(&cookbooks); err != nil {
		return nil, err
	}
	return cookbooks, nil
}

func (m *Cookbook) GetCookbook(id string) (*lib.Cookbook, error) {
	fmt.Println(id)
	var cookbook lib.Cookbook
	if err := m.Connect(); err != nil {
		return nil, err
	}
	if err := m.c.FindId(bson.ObjectIdHex(id)).One(&cookbook); err != nil {
		return nil, err
	}
	return &cookbook, nil
}

func (m *Cookbook) UpdateCookbook(id string, cookbook *lib.Cookbook) error {
	if err := m.Connect(); err != nil {
		return err
	}
	return m.c.UpdateId(bson.ObjectIdHex(id), cookbook)
}

func (m *Cookbook) DeleteCookbook(id string) error {
	if err := m.Connect(); err != nil {
		return err
	}
	return m.c.RemoveId(bson.ObjectIdHex(id))
}

func (m *Cookbook) Connect() (err error) {
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

func (m *Cookbook) Disconnect() {
	if m.session != nil {
		m.session.Close()
		m.session = nil
	}
}
