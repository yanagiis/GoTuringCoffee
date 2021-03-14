package entity

import (
	"GoTuringCoffee/internal/service/lib"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProcessBson Process mongo data structure
type ProcessBson struct {
	BaseItemBson
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Impl bson.Raw           `bson:"process_impl"`
}

// ToLibModel Convert process db model(ProcessBson) to lib.process(Process)
func (pb *ProcessBson) ToLibModel() (p lib.Process, err error) {
	processImpl, err := lib.NewProcessImpl(pb.Name)
	if err != nil {
		return lib.Process{}, err
	}

	b, err := pb.BaseItemBson.ToLibModel()
	if err != nil {
		log.Fatal("Failed to convert process BaseItemBson to lib BaseItem")
		return lib.Process{}, err
	}
	b.ID = ObjectIDToStringID(pb.ID)

	// Convert Bson.Raw to binary
	bson.Unmarshal(pb.Impl, processImpl)
	return lib.Process{
		BaseItem: b,
		Impl:     processImpl,
	}, nil
}

// CreateBsonFromProcessLibModel Convert lib process to process bson for mongodb
func CreateBsonFromProcessLibModel(p lib.Process) (pb ProcessBson, err error) {
	bb, err := CreateBaseItemFromLibModel(p.BaseItem)
	if err != nil {
		log.Fatalf("Falied to convert lib process %s to bson", p.Name)
		return
	}

	pb.BaseItemBson = bb
	pb.Impl, err = bson.Marshal(p.Impl)
	if err != nil {
		return
	}

	return
}
