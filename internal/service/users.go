package service

import (
	"context"
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

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repository      repository.UserRepository
	tokenManager    auth.TokenManager
	hasher          hash.PasswordHasher
	logger          *logging.Logger
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(repository repository.UserRepository, tokenManager auth.TokenManager, hasher hash.PasswordHasher,
	accessTokenTTL, refreshTokenTTL time.Duration, logger *logging.Logger) *UserService {
	return &UserService{
		repository:      repository,
		tokenManager:    tokenManager,
		hasher:          hasher,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		logger:          logger,
	}
}

func (s *UserService) Create(ctx context.Context, userDTO dto.CreateUserDTO) (dto.TokenDTO, error) {

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

func (s *UserService) FindOne(ctx context.Context, id string) (domain.User, error) {

	oid, err := params.ParseIdToObjectID(id)
	if err != nil {
		return domain.User{}, err
	}

	return s.repository.FindOne(ctx, oid)
}

func (s *UserService) FindAll(ctx context.Context, limit, offset, filter, sortBy string) (u []domain.User, err error) {
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

func (s *UserService) Update(ctx context.Context, userDTO dto.UpdateUserDTO) error {

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

func (s *UserService) Delete(ctx context.Context, id string) error {
	oid, err := params.ParseIdToObjectID(id)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, oid)
}

func (s *UserService) RefreshUserToken(ctx context.Context, userid string) (dto.TokenDTO, error) {

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

func (s *UserService) CreateSession(ctx context.Context, oid primitive.ObjectID) (dto.TokenDTO, error) {
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
