package app

import (
	"errors"
	"net/http"
	"test/internal/config"
	v1 "test/internal/delivery/http/v1"
	"test/internal/repository"
	"test/internal/server"
	"test/internal/service"
	"test/pkg/api/auth"
	"test/pkg/client/mongodb"
	"test/pkg/hash"
)

func Run() {

	cfg := config.GetConfig()

	mongoClient, err := mongodb.NewClient(cfg.MongodbConfig)
	if err != nil {
		return
	}

	db := mongoClient.Database(cfg.MongodbConfig.Database)
	repository := repository.NewRepository(db)

	tokenManager, err := auth.NewManager(cfg.AuthConfig.JWT.SecretKey)
	if err != nil {
		return
	}

	hasher := hash.NewSHA1Hasher(cfg.AuthConfig.PasswordSalt)

	services := service.NewServices(service.Deps{
		Repos:           repository,
		TokenManager:    tokenManager,
		Hasher:          hasher,
		AccessTokenTTL:  cfg.AuthConfig.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.AuthConfig.JWT.RefreshTokenTTL,
	})

	handlers := v1.NewHandler(services, tokenManager)

	router := handlers.Init()

	srv := server.NewServer(router, cfg)
	if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

}
