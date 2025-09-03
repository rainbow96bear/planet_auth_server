package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/rainbow96bear/planet_auth_server/oauth"
	"github.com/rainbow96bear/planet_db_server/logger"
	pb "github.com/rainbow96bear/planet_proto"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	providers map[string]oauth.Provider
}

func NewUserService(providers map[string]oauth.Provider) *UserService {
	return &UserService{providers: providers}
}

func (s *UserService) RefreshToken(ctx context.Context, req *pb.Token) (*pb.RefreshTokenResponse, error) {
	logger.Infof("start to refreshToken")
	logger.Debugf("receive token info : %+v", req)
	defer logger.Infof("end to refreshToken")

	parts := strings.SplitN(req.Token, "-", 2)
	if len(parts) < 2 {
		logger.Warnf("invalid token format: %s", req.Token)
		return nil, fmt.Errorf("invalid token format")
	}

	platform := parts[0]

	resRefreshToken := &pb.RefreshTokenResponse{}
	var err error
	switch platform {
	case "kakao":
		// 카카오 토큰 처리
		logger.Infof("handle kakao refresh token: %s", req.Token)
		resRefreshToken, err = s.providers["kakao"].RefreshToken(req.Token)
		if err != nil {
			return nil, err
		}
	case "google":
		// 구글 토큰 처리
		logger.Infof("handle google refresh token: %s", req.Token)
	default:
		logger.Warnf("unsupported platform: %s", platform)
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}

	// 처리 후 결과 반환
	return resRefreshToken, nil
}
