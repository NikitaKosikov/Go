package server

import (
	"net/http"
	"test/internal/config"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	writeTimeExpiration = 15 * time.Second
	readTimeExpiration  = 15 * time.Second
)

type server struct {
	httpServer *http.Server
}

func NewServer(router *gin.Engine, cfg *config.Config) *server {
	return &server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.ListenConfig.Port,
			Handler:      router,
			ReadTimeout:  readTimeExpiration,
			WriteTimeout: writeTimeExpiration,
		},
	}
}

func (s *server) Run() error {
	return s.httpServer.ListenAndServe()
}
