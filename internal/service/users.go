package service

import (
	"fmt"
	"test/internal/domain"
	"test/internal/repository"
	"test/internal/service/dto"
	"test/pkg/api"
	"test/pkg/api/auth"
	"test/pkg/api/params"
	"test/pkg/hash"
	"test/pkg/logging"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userService struct {
	repository      repository.UserRepository
	tokenManager    auth.TokenManager
	hasher          hash.PasswordHasher
	logger          *logging.Logger
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(repository repository.UserRepository, tokenManager auth.TokenManager, hasher hash.PasswordHasher,
	accessTokenTTL, refreshTokenTTL time.Duration, logger *logging.Logger) *userService {
	return &userService{
		repository:      repository,
		tokenManager:    tokenManager,
		hasher:          hasher,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		logger:          logger,
	}
}

func (s *userService) Create(ctx *gin.Context, userDTO dto.CreateUserDTO) (dto.TokenDTO, error) {

	if !dto.ValidCreateUserDTO(userDTO) {
		return dto.TokenDTO{}, fmt.Errorf("Invalid userDTO parameters")
	}

	passwordHash, err := s.hasher.Hash(userDTO.Password)
	if err != nil {
		return dto.TokenDTO{}, err
	}
	userDTO.Password = string(passwordHash)
	user := dto.ConvertCreateUserDTO(userDTO)
	id, err := s.repository.Create(ctx, user)
	if err != nil {
		return dto.TokenDTO{}, err
	}

	user.Id = id

	return s.CreateSession(ctx, user.Id)
}

func (s *userService) FindOne(ctx *gin.Context, id string) (domain.User, error) {

	oid, err := params.ParseIdToObjectID(id)
	if err != nil {
		return domain.User{}, err
	}

	return s.repository.FindOne(ctx, oid)
}

func (s *userService) FindAll(ctx *gin.Context, limit, offset, filter, sortBy string) (u []domain.User, err error) {
	filters, err := api.ParseFilters(filter)
	if err != nil {
		return u, err
	}
	for _, filter := range filters {
		if !ValidateUserField(filter.Field) {
			return u, fmt.Errorf("unknown field in sortBy query parameter")
		}
	}

	sortOptions, err := api.ParseSort(sortBy)
	if err != nil {
		return u, err
	}

	for _, option := range sortOptions {
		if !ValidateUserField(option.Field) {
			return u, fmt.Errorf("unknown field in sortBy query parameter")
		}
	}

	pagination, err := api.NewPagination(limit, offset)
	if err != nil {
		return u, err
	}
	return s.repository.FindAll(ctx, pagination, filters, sortOptions)
}

func (s *userService) Update(ctx *gin.Context, userDTO dto.UpdateUserDTO) error {

	if !dto.ValidUpdateUserDTO(userDTO) {
		return fmt.Errorf("Invalid userDTO parameters")
	}

	passwordHash, err := s.hasher.Hash(userDTO.Password)
	if err != nil {
		return err
	}
	userDTO.Password = string(passwordHash)
	user, err := dto.ConvertUpdateUserDTO(userDTO)
	if err != nil {
		return err
	}

	return s.repository.Update(ctx, user)
}

func (s *userService) Delete(ctx *gin.Context, id string) error {
	oid, err := params.ParseIdToObjectID(id)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, oid)
}

func (s *userService) RefreshUserToken(ctx *gin.Context, userid string) (dto.TokenDTO, error) {

	oid, err := params.ParseIdToObjectID(userid)
	if err != nil {
		return dto.TokenDTO{}, err
	}
	user, err := s.repository.GetUserByRefreshToken(ctx, oid)
	if err != nil {
		return dto.TokenDTO{}, err
	}
	return s.CreateSession(ctx, user.Id)
}

func (s *userService) CreateSession(ctx *gin.Context, oid primitive.ObjectID) (dto.TokenDTO, error) {
	id := oid.Hex()
	accessToken, err := s.tokenManager.GenerateAccessToken(id, s.accessTokenTTL)
	if err != nil {
		return dto.TokenDTO{}, err
	}
	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return dto.TokenDTO{}, err
	}

	session := domain.Session{
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	if err := s.repository.SetSession(ctx, oid, session); err != nil {
		return dto.TokenDTO{}, err
	}

	return dto.TokenDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
