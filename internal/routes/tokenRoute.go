package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/internal/handler"
)

func RegisterTokenRoutes(r *gin.Engine, tokenHandler *handler.TokenHandler) {
	oauthGroup := r.Group("/auth/token")
	oauthGroup.POST("/reissue/refresh", tokenHandler.ReissueRefreshToken)
	oauthGroup.POST("/issue/access", tokenHandler.IssueAccessToken)
}
