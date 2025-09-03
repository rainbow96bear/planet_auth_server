package kakao

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/grpc_client"
	"github.com/rainbow96bear/planet_auth_server/logger"
	"github.com/rainbow96bear/planet_auth_server/oauth/token"
	pb "github.com/rainbow96bear/planet_proto"
)

type Provider struct {
	RestApiKey   string
	RedirectUrl  string
	ClientSecret string
	AccessToken  map[string]string
}

type KakaoTokenResponse struct {
	AccessToken  string  `json:"access_token"`
	ExpiresIn    float64 `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"`
	Scope        string  `json:"scope"`
}

// 최상위 사용자 정보
type KakaoUser struct {
	ConnectedAt  time.Time     `json:"connected_at"`
	ID           int64         `json:"id"`
	KakaoAccount KakaoAccount  `json:"kakao_account"`
	Properties   KakaoProperty `json:"properties"`
}

// kakao_account 객체
type KakaoAccount struct {
	Profile                       KakaoProfile `json:"profile"`
	ProfileImageNeedsAgreement    bool         `json:"profile_image_needs_agreement"`
	ProfileNicknameNeedsAgreement bool         `json:"profile_nickname_needs_agreement"`
}

// profile 객체
type KakaoProfile struct {
	IsDefaultImage    bool   `json:"is_default_image"`
	IsDefaultNickname bool   `json:"is_default_nickname"`
	Nickname          string `json:"nickname"`
	ProfileImageURL   string `json:"profile_image_url"`
	ThumbnailImageURL string `json:"thumbnail_image_url"`
}

// properties 객체
type KakaoProperty struct {
	Nickname       string `json:"nickname"`
	ProfileImage   string `json:"profile_image"`
	ThumbnailImage string `json:"thumbnail_image"`
}

func (k *Provider) Signup(c *gin.Context) {
	logger.Infof("start get kakao signup")
	defer logger.Infof("end get kakao signup")

	code := c.Query("code")
	logger.Debugf("authorize code : %+v", code)
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code parameter is missing"})
		return
	}

	// access token 요청
	accessTokenUrl := fmt.Sprintf(
		"https://kauth.kakao.com/oauth/token?grant_type=authorization_code&client_id=%s&redirect_uri=%s&code=%s&client_secret=%s",
		k.RestApiKey,
		k.RedirectUrl,
		code,
		k.ClientSecret,
	)
	resp, err := http.Post(accessTokenUrl, "application/x-www-form-urlencoded", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to request access token"})
		return
	}
	defer resp.Body.Close()

	var tokenResp KakaoTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse access token response"})
		return
	}
	logger.Debugf("access token result : %+v", tokenResp)

	// access token으로 user info 요청
	req, err := http.NewRequest(
		http.MethodGet,
		"https://kapi.kakao.com/v2/user/me",
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	client := &http.Client{}
	userInfoResp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to request user info"})
		return
	}
	defer userInfoResp.Body.Close()

	var userInfoResult KakaoUser
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfoResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user info response"})
		return
	}

	logger.Debugf("user info result : %+v", userInfoResult)

	newUserInfo := &pb.UserInfo{
		OauthPlatform: "kakao",
		OauthId:       strconv.FormatInt(userInfoResult.ID, 10),
		Role:          "user",
		Nickname:      "행성",
		ProfileImage:  "http://localhost:5173/profile/image/default",
	}

	if userInfoResult.KakaoAccount.ProfileNicknameNeedsAgreement && userInfoResult.Properties.Nickname != "" {
		newUserInfo.Nickname = userInfoResult.Properties.Nickname
	}

	if userInfoResult.KakaoAccount.ProfileImageNeedsAgreement && userInfoResult.Properties.ProfileImage != "" {
		newUserInfo.ProfileImage = userInfoResult.Properties.ProfileImage
	}

	dbReq := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	logger.Debugf("send to db server about user info : %+v", newUserInfo)

	responseSignup, err := grpc_client.ReqOauthSignUp(dbReq, newUserInfo)
	if err != nil {
		logger.Warnf("failed to oauth sign up Error : %+v", err)
	}

	accessToken, err := token.CreateAccessToken(responseSignup.UserInfo.Id)
	if err != nil {

	}

	refreshToken, err := token.CreateRefreshToken()
	if err != nil {

	}
	// TODO : return 정의

	reqRefreshToken := &pb.Token{
		UserId: responseSignup.UserInfo.Id,
		Token:  refreshToken,
		Expiry: time.Now().Add(3 * 24 * time.Hour).Unix(),
	}
	logger.Debugf("reqRefreshToken value : %+v", reqRefreshToken)
	_, err = grpc_client.ReqRefreshToken(dbReq, reqRefreshToken)
	if err != nil {
		logger.Warnf("failed to refresh Token : %+v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    900,
	})

	redirectUrl := fmt.Sprintf("%s/login/callback", config.PLANET_CLIENT_ADDR)
	c.Redirect(http.StatusFound, redirectUrl)
}

func (k *Provider) Logout(c *gin.Context) {
	logger.Infof("start kakao logout")
	defer logger.Infof("end kakao logout")

	var req token.DeleteTokenRequest

	// JSON 바디 파싱
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("failed to parsing to JSON req : %+v", req)
		return
	}

	// 값 검증
	if req.UserID == 0 || req.Token == "" {
		logger.Warnf("user id or token is wrong userid : %v, token : %v", req.UserID, req.Token)
		return
	}

	if err := token.DeleteRefreshToken(req.UserID, req.Token); err != nil {
		logger.Warnf("failed to delete refresh token")
	}

	// TODO : redirect to home
	// TODO : refresh token에 따라 어떤 응답을 전달할지
}

func (k *Provider) RefreshToken(refreshToken string) (*pb.RefreshTokenResponse, error) {
	logger.Infof("start refresh token")
	defer logger.Infof("end refresh token")

	dbReq := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	logger.Debugf("request to db server for get user id : %+v", refreshToken)

	reqRefreshToken := &pb.Token{
		Token: refreshToken,
	}

	resRefreshTokenInfo, err := grpc_client.ReqGetRefreshTokenInfo(dbReq, reqRefreshToken)
	if err != nil {
		logger.Warnf("failed to oauth sign up Error : %+v", err)
	}

	accessTokenStr, err := token.CreateAccessToken(resRefreshTokenInfo.UserId)
	if err != nil {
		logger.Warnf("fail to create access token ERROR : %s", err.Error())
		return nil, err
	}

	newRefreshTokenStr, err := token.CreateRefreshToken()
	if err != nil {
		logger.Warnf("fail to create refresh token ERROR : %s", err.Error())
		return nil, err
	}

	resRefreshToken := &pb.RefreshTokenResponse{
		AccessToken:  accessTokenStr,
		RefreshToken: newRefreshTokenStr,
		Expiry:       time.Now().Add(time.Duration(config.ACCESS_TOKEN_EXPIRY_MINUTE) * time.Minute).Unix(),
	}

	logger.Debugf("reqRefreshToken value : %+v", reqRefreshToken)
	_, err = grpc_client.ReqRefreshToken(dbReq, reqRefreshToken)
	if err != nil {
		logger.Warnf("failed to refresh Token : %+v", err)
		return nil, err
	}

	return resRefreshToken, nil
}
