package oauth

import pb "github.com/rainbow96bear/planet_proto"

// 모든 플랫폼 provider가 따라야 하는 공통 인터페이스
type Provider interface {
	RefreshToken(token string) (*pb.RefreshTokenResponse, error)
}
