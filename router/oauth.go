package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/oauth/kakao"
)

func RegisterKakaoOauthRoutes(r *gin.Engine, provider *kakao.Provider) {

	oauthGroup := r.Group("/oauth/kakao")
	oauthGroup.GET("/signup", provider.Signup)
	oauthGroup.POST("/logout", provider.Logout)
}
