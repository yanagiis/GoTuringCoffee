package entity

import (
	"GoTuringCoffee/internal/service/lib"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseItemBson BaseItem bson version for mongodb
type BaseItemBson struct {
	Name        string   `bson:"name"`
	Description string   `bson:"description"`
	Tags        []string `bson:"tags,omitempty"`
	Notes       []string `bson:"notes,omitempty"`
	CreatedAt   int64    `bson:"created_at"`
	UpdatedAt   int64    `bson:"updated_at"`
}

// ToLibModel convert bson object to lib object
func (bb *BaseItemBson) ToLibModel() (b lib.BaseItem, err error) {
	b.Name = bb.Name
	b.Description = bb.Description
	b.Tags = bb.Tags
	b.Notes = bb.Notes
	b.CreatedAt = time.Unix(bb.CreatedAt, 0)
	b.UpdatedAt = time.Unix(bb.UpdatedAt, 0)

	return
}

// SetUpdatedTime update the updated time to now
func (bb *BaseItemBson) SetUpdatedTime() {
	bb.UpdatedAt = time.Now().UTC().Unix()
}

// CreateBaseItemFromLibModel Convert lib model to bson model
func CreateBaseItemFromLibModel(b lib.BaseItem) (bb BaseItemBson, err error) {
	bb.Name = b.Name
	bb.Description = b.Description
	bb.Tags = b.Tags
	bb.Notes = b.Notes
	bb.CreatedAt = b.CreatedAt.Unix()
	bb.UpdatedAt = b.UpdatedAt.Unix()

	return
}

// StringIDToObjectID Convert id string to ObjectID
func StringIDToObjectID(sid string) (oid primitive.ObjectID, err error) {
	if sid == "" {
		oid = primitive.NewObjectID()
	} else {
		oid, err = primitive.ObjectIDFromHex(sid)
	}

	return
}

// ObjectIDToStringID convert ObjectID to string id
func ObjectIDToStringID(oid primitive.ObjectID) (sid string) {
	return oid.Hex()
}
