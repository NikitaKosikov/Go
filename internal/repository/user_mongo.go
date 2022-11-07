package repository

import (
	"context"
	"errors"
	"fmt"
	"test/internal/domain"
	"test/pkg/api"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	descOrderKey   = "desc"
	descMongoDbKey = "-1"
	ascMongoDbKey  = "1"
)

var _ UserRepository = &userRepository{}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) UserRepository {
	return &userRepository{
		collection: database.Collection(usersCollection),
	}
}

// Create implements user.Storage
func (d *userRepository) Create(ctx context.Context, user domain.User) (primitive.ObjectID, error) {

	result, err := d.collection.InsertOne(ctx, &user)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("failed to create user due to error: %v", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}

	return primitive.ObjectID{}, fmt.Errorf("failed to create user")
}

// Delete implements user.Storage
func (d *userRepository) Delete(ctx context.Context, oid primitive.ObjectID) error {

	filter := bson.M{"_id": oid}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user with oid=%s due to error: %v", oid, err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

// FindAll implements user.Storage
func (d *userRepository) FindAll(ctx context.Context, p api.Pagination, filters []api.Filters, sort []api.Options) (u []domain.User, err error) {
	options := setSorting(sort).SetSkip(p.Offset)
	if p.Limit != 0 {
		options.SetLimit(p.Limit)
	}

	filter := setFilters(filters)
	cursor, err := d.collection.Find(ctx, filter, options)
	if err != nil {
		return u, fmt.Errorf("failed to find all users due to error:=%v", err)
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor due to error: %v", err)

	}

	return u, nil
}

// FindOne implements user.Storage
func (d *userRepository) FindOne(ctx context.Context, oid primitive.ObjectID) (u domain.User, err error) {

	filter := bson.M{"_id": oid}
	result := d.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, domain.ErrUserNotFound
		}
		return u, fmt.Errorf("failed to find user by oid=%s, due to error:=%v", oid, result.Err())
	}

	if err := result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user by oid=%s, from DB due to error: %v", oid, err)
	}

	return u, nil
}

// Update implements user.Storage
func (d *userRepository) Update(ctx context.Context, user domain.User) error {

	updateQuery := bson.M{}
	updateQuery["email"] = user.Email
	updateQuery["password"] = user.PasswordHash

	filter := bson.M{"_id": user.Id}

	result, err := d.collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateQuery}})
	if err != nil {
		return fmt.Errorf("failed to exceute update user query due to error: %v", err)
	}
	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) SetSession(ctx context.Context, oid primitive.ObjectID, session domain.Session) error {
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}}

	if _, err := r.collection.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("Failed to store session")
	}
	return nil
}

func (r *userRepository) GetUserByRefreshToken(ctx context.Context, oid primitive.ObjectID) (domain.User, error) {

	var user domain.User
	filter := bson.M{
		"_id":               oid,
		"session.expiresat": bson.M{"$gt": time.Now()},
	}
	if err := r.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed to find user by refresh token")

	}
	return user, nil
}
