package service

type TokenService struct {
	AccessTokenExpiry 
	RefreshTokenName
	RefreshTokenExpiry
	JwtSecretKey string
}

func (s *TokenService)IssueAccessToken(refreshToken string) (accessToken, error){
	// db에 refreshtoken으로 useruuid 얻기
	
	accessClaims := jwt.MapClaims{
		"userUuid":  userUuid,
		"plateform": "kakao",
		"exp":       time.Now().Add(time.Duration(s.AccessTokenExpiry) * time.Minute),
		"iat":       time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	return accessTokenString, nil
}

func (s *TokenService)IssueRefreshToken(userUuid string) (string, error) {
	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)
	refreshToken = fmt.Sprintf("%s", refreshToken)

	// db에 refresh token 저장
	result, err := db.UpdateRefreshToken(userUuid, refreshToken)
    if err != nil {
        logger.Errorf("update refresh token ERR[%s]", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		redirectUrl := fmt.Sprintf("%s/login/callback?status=error&code=%s", config.PLANET_CLIENT_ADDR, utils.ERR_REFRESH_TOKEN_CREATE)
		c.Redirect(http.StatusFound, redirectUrl)
        return
    }
	return refreshToken
}