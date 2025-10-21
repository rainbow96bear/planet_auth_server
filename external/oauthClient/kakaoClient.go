package oauthClient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"planet_utils/pkg/logger"
)

type KakaoClient struct {
	RestApiKey   string
	RedirectUrl  string
	ClientSecret string
}

func (c *KakaoClient) GetAccessToken(code string) (string, error) {
	logger.Infof("start to get access token")
	defer logger.Infof("end to get access token")
	accessTokenUrl := fmt.Sprintf(
		"https://kauth.kakao.com/oauth/token?grant_type=authorization_code&client_id=%s&redirect_uri=%s&code=%s&client_secret=%s",
		c.RestApiKey,
		c.RedirectUrl,
		code,
		c.ClientSecret,
	)

	resp, err := http.Post(accessTokenUrl, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp KakaoTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}
	logger.Debugf("access token result : %+v", tokenResp)

	return tokenResp.AccessToken, nil
}

func (c *KakaoClient) GetUserInfo(accessToken string) (*KakaoUser, error) {
	logger.Infof("start to get user info")
	defer logger.Infof("end to get user info")
	req, err := http.NewRequest(
		http.MethodGet,
		"https://kapi.kakao.com/v2/user/me",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	userInfoResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer userInfoResp.Body.Close()

	var userInfoResult *KakaoUser
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfoResult); err != nil {
		return nil, err
	}

	return userInfoResult, nil
}
