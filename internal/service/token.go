package service

type TokenService struct {
	AccessTokenExpiry 
	RefreshTokenName
	RefreshTokenExpiry
	JwtSecretKey string
}

func (s *TokenService)IssueAccessToken(userUuid string) (accessToken, error){
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

func (s *TokenService)IssueRefreshToken() string {
	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)
	refreshToken = fmt.Sprintf("%s", refreshToken)
	return refreshToken
}