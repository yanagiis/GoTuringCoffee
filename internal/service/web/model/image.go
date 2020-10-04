package model

import "github.com/globalsign/mgo/bson"

type MongoImage struct {
	ID          bson.ObjectId `bson:"_id"`
	Author      string        `bson:"author"`
	Caption     string        `bson:"caption"`
	ContentType string        `bson:"contentType"`
	DateTime    string        `bson:"dateTime"`
	FileID      bson.ObjectId `bson:"fileID"`
	FileSize    int64         `bson:"fileSize"`
	Height      int           `bson:"height"`
	Name        string        `bson:"name"`
	Width       int           `bson:"width"`
}
