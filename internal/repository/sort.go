package repository

import (
	"strings"
	"test/pkg/api"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


const (
	descMongoDbKey = "-1"
	ascMongoDbKey  = "1"
	DescOrderKey   = "desc"
)

func setSorting(sortOptions []api.Options) *options.FindOptions {
	if sortOptions==nil || len(sortOptions)==0{
		return &options.FindOptions{}
	}
	var filterOptions bson.A
	for _, option := range sortOptions {
		order := convertSortOrder(option.Order)
		filterOptions = append(filterOptions, bson.D{{Key: option.Field, Value: order}})
	}
	return options.Find().SetSort(filterOptions)
}

func convertSortOrder(order string) string {
	if strings.ToLower(order) == DescOrderKey {
		return descMongoDbKey
	}
	return ascMongoDbKey
}
