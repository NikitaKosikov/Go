package v1

import (
	"test/internal/service"
	"test/pkg/api/auth"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.GET("/auth/google/login", OauthGoogleLogin)
	router.GET("/auth/google/callback", OauthGoogleCallback)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	api := router.Group(auth.BasicURL + auth.Version)
	{
		h.initUsersRoutes(api)
	}
}
