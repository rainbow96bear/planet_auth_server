package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/internal/handler"
)

func RegisterKakaoOauthRoutes(r *gin.Engine, kakaoHandler *handler.KakaoHandler) {
	oauthGroup := r.Group("/oauth/kakao")
	oauthGroup.GET("/login", kakaoHandler.Login)
	oauthGroup.POST("/logout", kakaoHandler.Logout)
}
