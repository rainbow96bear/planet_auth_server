package handler

import (
    "net/http"
    "planet/external/oauth"
    "planet/internal/auth/service"
)

type KakaoHandler struct {
    KakaoClient  *oauth.KakaoClient
    TokenService *service.TokenService
}

func (h *KakaoHandler) Login(c *gin.Context) {
    logger.Infof("start kakao login")
	defer logger.Infof("end kakao login")

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

    // db에 사용자 정보 조회
    // platform and platform_id
    // 조회된 내용이 없으면 아직 가입 안 한 사용자
    // 가입 안 한 사용자는 session을 생성하고 signup으로 redirect

    refreshToken := h.TokenService.IssueRefreshToken(userUuid)

    // db에 refresh token을 사용자 uuid와 함께 저장 요청


    c.SetCookie(
		config.REFRESH_TOKEN_NAME,
		refreshToken,
		config.REFRESH_TOKEN_EXPIRY_DURATION,
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

    refreshToken, err := c.Cookie(config.REFRESH_TOKEN_NAME)
	if err != nil {
		logger.Warnf("no refresh_token cookie found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

    // db에서 refresh token 삭제
    result, err := db.DeleteRefreshToken(refreshToken)
    if err != nil {
        logger.Warnf("delete refresh token ERR[%s]", err.Error())
    }

    c.SetCookie(config.REFRESH_TOKEN_NAME, "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}