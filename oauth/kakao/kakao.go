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

func (k *Provider) Logout(c *gin.Context) {
	// logger.Infof("start kakao logout")
	// defer logger.Infof("end kakao logout")
	// accessToken := k.AccessToken[userId] // 실제 발급받은 Access Token으로 교체

	// req, err := http.NewRequest(
	// 	http.MethodPost,
	// 	"https://kapi.kakao.com/v1/user/logout",
	// 	nil,
	// )
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
	// 	return
	// }

	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	// req.Header.Set("Authorization", "Bearer "+accessToken)

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to request logout"})
	// 	return
	// }
	// defer resp.Body.Close()

	// var result map[string]interface{}
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse logout response"})
	// 	return
	// }

	// c.JSON(http.StatusOK, result)
}

func (k *Provider) GetUserInfo(c *gin.Context) {
	logger.Infof("start get kakao user info")
	defer logger.Infof("end get kakao user info")

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

	_, err = grpc_client.ReqOauthSignUp(dbReq, newUserInfo)
	if err != nil {
		logger.Warnf("failed to oauth sign up Error : %+v", err)
	}

	c.JSON(http.StatusOK, userInfoResult)
}
