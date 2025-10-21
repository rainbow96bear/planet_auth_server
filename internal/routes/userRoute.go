package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/internal/handler"
)

func RegisterUserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	userGroup := r.Group("/auth/user")
	userGroup.POST("/signup", userHandler.Signup)
	userGroup.GET("/nickname/check", userHandler.NicknameCheck)
}
