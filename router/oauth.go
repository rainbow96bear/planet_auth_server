package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/oauth/kakao"
)

func RegisterKakaoOauthRoutes(r *gin.Engine, provider *kakao.OauthProvider) {
	oauthGroup := r.Group("/oauth/kakao")
	oauthGroup.GET("/login", provider.Login)
	oauthGroup.POST("/logout", provider.Logout)
}
