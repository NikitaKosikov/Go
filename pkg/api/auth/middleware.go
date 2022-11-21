package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	RefreshURL          = "/:id/auth/refresh"
	BasicURL            = "/api"
	Version             = "/v1"
	authorizationHeader = "Authorization"
	userURLAPI          = "http://localhost:4000/api/v1/users/:id"
	refreshTokenURI     = "/auth/refresh"
)

func (m *Manager) VerifyJWTMiddleware(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtToken, err := parseAuthHeader(ctx)
		if err != nil {
			return
		}
		userid := ctx.Param(IdNameURL)
		claims := &Claims{}
		token, err := m.GetTokenFromString(jwtToken, claims)
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, userURLAPI+userid+refreshTokenURI)
			return
		}
		if err := m.ValidateToken(token, claims); err != nil {
			ctx.Status(http.StatusForbidden)
			ctx.Writer.Write([]byte(err.Error()))
			ctx.Next()
		}
		if !hasPermission(roles, claims, userid) {
			ctx.Status(http.StatusForbidden)
			ctx.Writer.Write([]byte("Forbidden"))
			ctx.Next()
		}
		return
	}
	
}

func parseAuthHeader(ctx *gin.Context) (string, error) {
	if ctx.GetHeader(authorizationHeader) == "" {
		return "", fmt.Errorf("empty auth header")
	}
	authHeader := strings.Split(ctx.GetHeader(authorizationHeader), PrefixToken)
	if len(authHeader) != 2 {
		return "", fmt.Errorf("invalid auth header")
	}
	if len(authHeader[1]) == 0 {
		return "", fmt.Errorf("token is empty")
	}
	return authHeader[1], nil
}

func hasPermission(roles []string, claims *Claims, id string) bool {
	if claims.Id != id {
		return false
	}
	for _, role := range roles {
		if strings.ToLower(claims.Role) == strings.ToLower(role) {
			return true
		}
	}
	return false
}
