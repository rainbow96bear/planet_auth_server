package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/auth"
	"github.com/rainbow96bear/planet_auth_server/auth/token"
)

func RegisterSignupRoutes(r *gin.Engine) {
	r.POST("/auth/signup", auth.Signup)
	r.GET("/auth/nickname/available", auth.NicknameAvailable)
}

func RegisterTokenRoutes(r *gin.Engine, provider *token.TokenProvider) {
	tokenGroup := r.Group("/auth/token")
	tokenGroup.POST("/refresh", provider.UpdateRefreshTokens)
	tokenGroup.POST("/access", provider.IssueAccessToken)
}
