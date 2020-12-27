package entity

import (
	"GoTuringCoffee/internal/service/lib"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ProcessBson Process mongo data structure
type ProcessBson struct {
	ID        string   `bson:"id, omitempty"`
	Name      string   `bson:"name"`
	Impl      bson.Raw `bson:"process_impl"`
	CreatedAt int64    `bson:"created_at"`
	UpdatedAt int64    `bson:"updated_at"`
}

// ToLibModel Convert process db model(ProcessBson) to lib.process(Process)
func (pb *ProcessBson) ToLibModel() (p lib.Process, err error) {
	processImpl, err := lib.NewProcessImpl(pb.Name)
	if err != nil {
		return lib.Process{}, err
	}

	// Convert Bson.Raw to binary
	bson.Unmarshal(pb.Impl, processImpl)
	return lib.Process{
		ID:        pb.ID,
		Name:      pb.Name,
		CreatedAt: time.Unix(pb.CreatedAt, 0),
		UpdatedAt: time.Unix(pb.UpdatedAt, 0),
		Impl:      processImpl,
	}, nil
}
