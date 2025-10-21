package oauthClient

import "time"

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
