package repository

import (
	"test/internal/domain"
	"test/pkg/api"
	"test/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(c *gin.Context, user domain.User) (primitive.ObjectID, error)
	FindOne(c *gin.Context, oid primitive.ObjectID) (domain.User, error)
	FindAll(c *gin.Context, pagination api.Pagination, filters []api.Filters, sortOptions []api.Options) (u []domain.User, err error)
	Update(c *gin.Context, user domain.User) error
	Delete(c *gin.Context, oid primitive.ObjectID) error
	SetSession(c *gin.Context, oid primitive.ObjectID, session domain.Session) error
	GetUserByRefreshToken(c *gin.Context, id primitive.ObjectID) (domain.User, error)
}

type Repository struct {
	UserRepositiry UserRepository
}

func NewRepository(db *mongo.Database, logger *logging.Logger) *Repository {
	return &Repository{
		UserRepositiry: NewUserRepository(db, logger),
	}
}
