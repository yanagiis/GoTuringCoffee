package entity

import (
	"GoTuringCoffee/internal/service/lib"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GroupBson Group bson
type GroupBson struct {
	BaseItemBson
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	Cookbooks []primitive.ObjectID `bson:"cookbooks"`
	Groups    []primitive.ObjectID `bson:"groups"`
}

// ToLibModel Convert bson to lib model
func (gb GroupBson) ToLibModel() (g lib.Group, err error) {
	b, err := gb.BaseItemBson.ToLibModel()
	if err != nil {
		log.Fatalf("Failed to convert cookbook bson %s to lib cookbook", gb.Name)
		return
	}
	b.ID = ObjectIDToStringID(gb.ID)
	g.BaseItem = b

	// Convert all groups and cookbooks id to string
	for _, oid := range gb.Cookbooks {
		g.Cookbooks = append(g.Cookbooks, ObjectIDToStringID(oid))
	}

	for _, oid := range gb.Groups {
		g.Cookbooks = append(g.Groups, ObjectIDToStringID(oid))
	}

	return
}

// CreateBsonFromGroupLibModel Convert lib model to bson
func CreateBsonFromGroupLibModel(group lib.Group) (gb GroupBson, err error) {
	bb, err := CreateBaseItemFromLibModel(group.BaseItem)
	if err != nil {
		log.Fatalf("Failed to convert cookbook %s to bson model", group.Name)
		return
	}

	gb.ID, err = StringIDToObjectID(group.BaseItem.ID)
	if err != nil {
		log.Fatalf("Failed to convert string id %s to object id", group.ID)
		return
	}
	gb.BaseItemBson = bb

	// Convert all string id of groups and cookbooks to objectID
	for _, sid := range group.Groups {
		oid, err := StringIDToObjectID(sid)
		if err != nil {
			log.Fatalf("Failed to convert sid of group %s to objectID", sid)
		}
		gb.Cookbooks = append(gb.Cookbooks, oid)
	}

	for _, sid := range group.Cookbooks {
		oid, err := StringIDToObjectID(sid)
		if err != nil {
			log.Fatalf("Failed to convert sid of cookbook %s to objectID", sid)
		}
		gb.Cookbooks = append(gb.Cookbooks, oid)
	}

	return
}
