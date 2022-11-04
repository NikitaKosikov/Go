package service

import (
	"test/internal/domain"
	"test/internal/repository"
	"test/internal/service/dto"
	"test/pkg/api/auth"
	"test/pkg/hash"
	"test/pkg/logging"

	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Users interface {
	Create(ctx *gin.Context, userDTO dto.CreateUserDTO) (dto.TokenDTO, error)
	FindOne(ctx *gin.Context, id string) (domain.User, error)
	FindAll(ctx *gin.Context, limit, offset, filter, sortBy string) (u []domain.User, err error)
	Update(ctx *gin.Context, userDTO dto.UpdateUserDTO) error
	Delete(ctx *gin.Context, id string) error
	RefreshUserToken(ctx *gin.Context, userId string) (dto.TokenDTO, error)
	CreateSession(ctx *gin.Context, oid primitive.ObjectID) (dto.TokenDTO, error)
}

type Deps struct {
	Repos           *repository.Repository
	TokenManager    auth.TokenManager
	Hasher          hash.PasswordHasher
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Services struct {
	Users           Users
	hasher          hash.PasswordHasher
	AccessTokenTtl  time.Duration
	RefreshTokenTtl time.Duration
}

func NewServices(deps Deps, logger *logging.Logger) *Services {
	usersService := NewUserService(deps.Repos.UserRepositiry, deps.TokenManager, deps.Hasher,
		deps.AccessTokenTTL, deps.RefreshTokenTTL, logger)
	return &Services{
		Users: usersService,
	}
}
