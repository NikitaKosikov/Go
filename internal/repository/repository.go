package repository

import (
	"context"
	"test/internal/domain"
	"test/pkg/api"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go -package=mocks
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (primitive.ObjectID, error)
	FindOne(ctx context.Context, oid primitive.ObjectID) (domain.User, error)
	FindAll(ctx context.Context, pagination api.Pagination, filters []api.Filters, sortOptions []api.Options) (u []domain.User, err error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, oid primitive.ObjectID) error
	SetSession(ctx context.Context, oid primitive.ObjectID, session domain.Session) error
	GetUserByRefreshToken(ctx context.Context, id primitive.ObjectID) (domain.User, error)
}

type Repository struct {
	UserRepositiry UserRepository
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		UserRepositiry: NewUserRepository(db),
	}
}
