package entity

import (
	"GoTuringCoffee/internal/service/lib"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

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

// CreateBsonFromLibModel create bson from lib model
func CreateBsonFromLibModel(cookbook lib.Cookbook) (CookbookBson, error) {
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
