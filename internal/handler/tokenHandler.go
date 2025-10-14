package handler

import (
    "net/http"
    "planet/external/oauth"
    "planet/internal/auth/service"
)

type TokenHandler struct {
    TokenService *service.TokenService
}

func (h *TokenHandler) IssueAccessToken(c *gin.Context) {
    logger.Infof("start to issue access token")
	defer logger.Infof("end to issue access token")

	refreshToken, err := c.Cookie(config.REFRESH_TOKEN_NAME)
	if err != nil {
		logger.Warnf("no %s cookie found", config.REFRESH_TOKEN_NAME)
		return
	}

	accessToken, err := h.TokenService.IssueAccessToken(refreshToken)
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

func (h *TokenHandler) IssueRefreshToken(c *gin.Context) {
    logger.Infof("start to issue refresh token")
	defer logger.Infof("end to issue refresh token")

	refreshToken, err := c.Cookie(config.REFRESH_TOKEN_NAME)
	if err != nil {
		logger.Warnf("no refresh_token cookie found")
		return
	}

	accessToken, err := h.TokenService.IssueAccessToken(refreshToken)
	if err != nil {
		logger.Warnf("fail to create access token ERROR : %s", err.Error())
		return
	}

	refreshToken, err := h.TokenService.IssueRefreshToken()
	if err != nil {
		logger.Warnf("no refresh_token cookie found")
		return
	}
	
	c.SetCookie(
		h.TokenService.RefreshTokenName,
		refreshToken,
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