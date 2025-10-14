package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/oauth/kakao"
	"planet/internal/api/oauth/handler"
)

func RegisterTokenRoutes(r *gin.Engine, tokenHandler *handler.TokenHandler) {
	oauthGroup := r.Group("/auth/token")
	oauthGroup.POST("/issue/refresh", tokenHandler.IssueRefreshToken)
	oauthGroup.POST("/issue/access", tokenHandler.IssueAccessToken)
}