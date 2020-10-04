package model

import "go.mongodb.org/mongo-driver/bson"

// SortDescendingCreatedAt returns the configuration for having descending elements
var SortDescendingCreatedAt = bson.M{"created_at": -1}
