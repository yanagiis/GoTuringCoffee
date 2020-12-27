package lib

import (
	"testing"

	"github.com/globalsign/mgo"
)

func TestInitMgo(t *testing.T) {
	url := "mongodb+srv://turingcoffee:test12345@cluster0.m5idb.gcp.mongodb.net/testturingcoffee?retryWrites=true&w=majority"
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(database).C(collection)
	err := c.Find(query).One(&result)
	session, err := mgo.Dial(url)

}
