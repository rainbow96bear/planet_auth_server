package oauth

import (
	"github.com/rainbow96bear/planet_auth_server/utils"
)

// 모든 플랫폼 provider가 따라야 하는 공통 인터페이스
type OauthProvider interface {
	utils.Provider
	// RotateRefreshToken(token string) (*pb.RefreshTokenResponse, error)
}
