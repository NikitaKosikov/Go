package repository

import (
	"test/pkg/api"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setFilters(filters []api.Filters) primitive.D {
	if filters == nil || len(filters) == 0 {
		return primitive.D{}
	}
	var filter bson.D
	for _, f := range filters {
		filter = append(filter, bson.E{Key: f.Field, Value: bson.E{Key: "$" + f.Operation, Value: f.Value}})
	}
	return filter
}
