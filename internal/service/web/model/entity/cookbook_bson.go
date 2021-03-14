package entity

import (
	"GoTuringCoffee/internal/service/lib"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CookbookBson Cookbook data structure
type CookbookBson struct {
	BaseItemBson
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Processes []ProcessBson      `bson:"processes,omitempty"`
}

// CreateBsonFromCookbookLibModel create bson from lib model
func CreateBsonFromCookbookLibModel(cookbook lib.Cookbook) (cb CookbookBson, err error) {
	bb, err := CreateBaseItemFromLibModel(cookbook.BaseItem)
	if err != nil {
		log.Fatalf("Failed to convert cookbook %s to bson model", cookbook.Name)
		return
	}

	cb.ID, err = StringIDToObjectID(cookbook.BaseItem.ID)
	if err != nil {
		log.Fatalf("Failed to convert string id %s to object id", cookbook.ID)
		return
	}
	cb.BaseItemBson = bb

	var pb ProcessBson
	for _, p := range cookbook.Processes {

		pb, err = CreateBsonFromProcessLibModel(p)
		if err != nil {
			return
		}
		pb.ID, err = StringIDToObjectID(cookbook.ID)

		cb.Processes = append(cb.Processes, pb)
	}

	return cb, nil
}

// ToLibModel Convert db model(CookbookBson) to lib.Cookbook(Cookbook)
func (cb CookbookBson) ToLibModel() (c lib.Cookbook, err error) {
	b, err := cb.BaseItemBson.ToLibModel()
	if err != nil {
		log.Fatalf("Failed to convert cookbook bson %s to lib cookbook", cb.Name)
		return
	}
	b.ID = ObjectIDToStringID(cb.ID)
	c.BaseItem = b

	for _, pj := range cb.Processes {
		p, err := pj.ToLibModel()
		if err != nil {
			continue
		}
		c.Processes = append(c.Processes, p)
	}
	return
}
