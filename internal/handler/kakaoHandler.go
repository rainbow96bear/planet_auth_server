package handler

import (
	"fmt"
	"net/http"
	"planet_utils/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/dto"
	"github.com/rainbow96bear/planet_auth_server/external/oauthClient"
	"github.com/rainbow96bear/planet_auth_server/internal/service"
)

type KakaoHandler struct {
	KakaoClient  *oauthClient.KakaoClient
	UserService  *service.UserService
	TokenService *service.TokenService
	Platform     string
}

func (h *KakaoHandler) Login(c *gin.Context) {
	logger.Infof("start kakao login")
	defer logger.Infof("end kakao login")
	ctx := c.Request.Context()
	code := c.Query("code")
	logger.Debugf("authorize code : %+v", code)
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code parameter is missing"})
		return
	}

	accessToken, err := h.KakaoClient.GetAccessToken(code)
	if err != nil {
		logger.Errorf("get access token ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userInfo, err := h.KakaoClient.GetUserInfo(accessToken)
	if err != nil {
		logger.Errorf("get user info ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userIDStr := strconv.FormatInt(userInfo.ID, 10)

	oauthUserInfo := &dto.OauthUserInfo{
		OauthPlatform: h.Platform,
		OauthId:       userIDStr,
	}

	userUuid, err := h.UserService.IsUserExists(ctx, oauthUserInfo)
	if err != nil {
		logger.Errorf("fail to check user info ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if userUuid == "" {
		oauthSession, err := h.UserService.CreateOauthSession(ctx, oauthUserInfo)
		if err != nil {
			logger.Warnf("failed to save oauth session: %+v", err)
			redirectUrl := fmt.Sprintf("%s/login/callback?status=error", config.PLANET_CLIENT_ADDR)
			c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
			return
		}
		c.SetCookie(
			"signup_session",
			oauthSession,
			1800, // 30분
			"/",
			"",
			true, // Secure
			true, // HttpOnly
		)

		// 3) 회원가입 페이지로 리다이렉트
		redirectUrl := fmt.Sprintf("%s/signup", config.PLANET_CLIENT_ADDR)
		c.Redirect(http.StatusFound, redirectUrl)
		return
	}

	refreshToken, err := h.TokenService.IssueRefreshToken(ctx, userUuid)

	c.SetCookie(
		config.REFRESH_TOKEN_NAME,
		refreshToken,
		int(config.REFRESH_TOKEN_EXPIRY_DURATION),
		"/",
		"",
		true,
		true,
	)

	redirectUrl := fmt.Sprintf("%s/login/callback?status=success", config.PLANET_CLIENT_ADDR)
	c.Redirect(http.StatusFound, redirectUrl)
}

func (h *KakaoHandler) Logout(c *gin.Context) {
	logger.Infof("start kakao logout")
	defer logger.Infof("end kakao logout")
	ctx := c.Request.Context()
	refreshToken, err := c.Cookie(config.REFRESH_TOKEN_NAME)
	if err != nil {
		logger.Warnf("no refresh_token cookie found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	err = h.TokenService.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		logger.Warnf("fail to revoke refresh token : %s", refreshToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "fail to revoke refresh token"})
		return
	}

	c.SetCookie(config.REFRESH_TOKEN_NAME, "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
