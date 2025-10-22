package handler

import (
	"net/http"
	"planet_utils/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/internal/service"
)

type TokenHandler struct {
	TokenService *service.TokenService
}

func (h *TokenHandler) IssueAccessToken(c *gin.Context) {
	logger.Infof("start to issue access token")
	defer logger.Infof("end to issue access token")
	ctx := c.Request.Context()
	refreshToken, err := c.Cookie(config.REFRESH_TOKEN_NAME)
	if err != nil {
		logger.Warnf("no %s cookie found", config.REFRESH_TOKEN_NAME)
		return
	}
	accessToken, err := h.TokenService.IssueAccessToken(ctx, refreshToken)
	if err != nil {
		logger.Warnf("fail to create access token ERR[%s]", err.Error())
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		h.TokenService.AccessTokenExpiry, // 15분
		"/",
		"",   // domain
		true, // secure
		true, // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   h.TokenService.AccessTokenExpiry,
	})
}

func (h *TokenHandler) ReissueRefreshToken(c *gin.Context) {
	logger.Infof("start to issue refresh token")
	defer logger.Infof("end to issue refresh token")
	ctx := c.Request.Context()
	refreshToken, err := c.Cookie(config.REFRESH_TOKEN_NAME)
	if err != nil {
		logger.Warnf("no refresh_token cookie found")
		return
	}

	accessToken, err := h.TokenService.IssueAccessToken(ctx, refreshToken)
	if err != nil {
		logger.Warnf("fail to create access token ERROR : %s", err.Error())
		return
	}

	newRefreshToken, err := h.TokenService.ReissueRefreshToken(ctx, refreshToken)
	if err != nil {
		logger.Warnf("no refresh_token cookie found")
		return
	}

	c.SetCookie(
		h.TokenService.RefreshTokenName,
		newRefreshToken,
		h.TokenService.RefreshTokenExpiry,
		"/",
		"",
		true,
		true,
	)

	c.SetCookie(
		"access_token",
		accessToken,
		h.TokenService.AccessTokenExpiry, // 15분
		"/",
		"",   // domain
		true, // secure
		true, // httpOnly
	)
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   h.TokenService.AccessTokenExpiry,
	})
}
