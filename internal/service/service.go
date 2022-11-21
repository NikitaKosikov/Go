package service

import (
	"context"
	"test/internal/domain"
	"test/internal/repository"
	"test/internal/service/dto"
	"test/pkg/api/auth"
	"test/pkg/hash"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go -package=mocks
type Users interface {
	Create(ctx context.Context, userDTO dto.CreateUserDTO) (dto.TokenDTO, error)
	FindOne(ctx context.Context, id string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindAll(ctx context.Context, limit, offset, filter, sortBy string) (u []domain.User, err error)
	Update(ctx context.Context, userDTO dto.UpdateUserDTO) error
	Delete(ctx context.Context, id string) error
	RefreshUserToken(ctx context.Context, userId string) (dto.TokenDTO, error)
	CreateSession(ctx context.Context, oid primitive.ObjectID) (dto.TokenDTO, error)
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
}

func NewServices(deps Deps) *Services {
	usersService := NewUserService(deps.Repos.UserRepositiry, deps.TokenManager, deps.Hasher,
		deps.AccessTokenTTL, deps.RefreshTokenTTL)
	return &Services{
		Users: usersService,
	}
}
