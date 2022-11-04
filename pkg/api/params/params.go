package params

import (
	"errors"
	"math"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	maxParamId = math.MaxInt32
)

var ErrInvalidIdParam = errors.New("invalid id param")

func ParseIdToObjectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, ErrInvalidIdParam
	}
	return oid, nil
}
