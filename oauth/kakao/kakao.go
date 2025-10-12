package kakao

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/auth/token"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/grpc_client"
	"github.com/rainbow96bear/planet_auth_server/logger"
	"github.com/rainbow96bear/planet_auth_server/utils"
	pb "github.com/rainbow96bear/planet_proto"
)

const PLATFORM_NAME = "kakao"

func (k *OauthProvider) Login(c *gin.Context) {
	logger.Infof("start get kakao login")
	defer logger.Infof("end get kakao login")

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

	// 사용자 정보 획득

	// oauthPlatform과 OauthId로 사용자 조회
	platformInfo := &pb.PlatformInfo{
		Platform:   PLATFORM_NAME,
		PlatformId: strconv.FormatInt(userInfoResult.ID, 10),
	}

	dbClient, err := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	if err != nil {
		logger.Errorf("fail to create dbClient ERR[%s]", err.Error())
		return
	}
	// 사용자 정보가 없으면 일단 생성 후 응답
	responseUserInfo, err := dbClient.ReqGetUserInfoByPlatformInfo(platformInfo)
	if err != nil {
		logger.Warnf("failed to get user info Error : %+v", err)
	}
	logger.Debugf("response user info by platform info : %+v", responseUserInfo)
	// 닉네임이 없는 사용자 == 처음 가입
	if responseUserInfo.Nickname == "" {
		// 1) 회원가입 필요 → 세션 생성
		// user uuid를 만들어서 users에도 저장 + oauth session 저장
		sessionID := utils.GenerateRandomSessionID()
		oauthSession := &pb.OauthSession{
			SessionId:  sessionID,
			Platform:   PLATFORM_NAME,
			PlatformId: strconv.FormatInt(userInfoResult.ID, 10),
		}

		responseOauthSession, err := dbClient.ReqSaveOauthSession(oauthSession)
		if err != nil || !responseOauthSession.Status {
			logger.Warnf("failed to save oauth session: %+v", err)
			// gRPC 실패 시 안전하게 에러 페이지나 기본 로그인 페이지로 redirect
			redirectUrl := fmt.Sprintf("%s/login/callback?status=error&code=%s", config.PLANET_CLIENT_ADDR, utils.ERR_DB_REQUEST)
			c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
			return
		}

		// 2) 쿠키 발급
		c.SetCookie(
			"signup_session",
			sessionID,
			1800, // 30분
			"/",
			"",
			true, // Secure
			true, // HttpOnly
		)

		// 3) 회원가입 페이지로 리다이렉트
		redirectUrl := fmt.Sprintf("%s/signup", config.PLANET_CLIENT_ADDR)
		c.Redirect(http.StatusFound, redirectUrl)
		return // 여기서 종료
	}

	// 4) 이미 가입된 경우 → 토큰 발급 후 로그인 완료 리다이렉트
	refreshToken := token.CreateRefreshToken()
	if err != nil {
		logger.Warnf("fail to issue tokens ERR[%s]", err.Error())
	}
	reqRefreshToken := &pb.Token{
		UserUuid: responseUserInfo.UserUuid,
		Token:    refreshToken,
		Expiry:   time.Now().Add(3 * 24 * time.Hour).Unix(),
	}

	_, err = dbClient.ReqUpdateRefreshToken(reqRefreshToken)
	if err != nil {
		logger.Warnf("failed to refresh Token : %+v", err)
		redirectUrl := fmt.Sprintf("%s/login/callback?status=error&code=%s", config.PLANET_CLIENT_ADDR, utils.ERR_REFRESH_TOKEN_CREATE)
		c.Redirect(http.StatusFound, redirectUrl)

	}

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

func (k *OauthProvider) Logout(c *gin.Context) {
	logger.Infof("start kakao logout")
	defer logger.Infof("end kakao logout")

	refreshToken, err := c.Cookie(config.REFRESH_TOKEN_NAME)
	if err != nil {
		logger.Warnf("no refresh_token cookie found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	// DB에서 삭제
	if err := token.DeleteRefreshToken(refreshToken); err != nil {
		logger.Warnf("failed to delete refresh token: %v", err)
	}

	// 쿠키 무효화
	c.SetCookie(config.REFRESH_TOKEN_NAME, "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
