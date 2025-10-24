package service

import (
	"context"
	"errors"
	"time"

	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/dto"
	"github.com/rainbow96bear/planet_auth_server/internal/repository"
	"github.com/rainbow96bear/planet_auth_server/utils"
	"github.com/rainbow96bear/planet_utils/model"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type UserService struct {
	// AccessTokenExpiry
	// JwtSecretKey string
	ProfileImgSavePath string
	UsersRepo          *repository.UsersRepository
	OauthSessionsRepo  *repository.OauthSessionsRepository
}

func (s *UserService) Signup(ctx context.Context, oauthSessionUuid string, signupInfo *dto.SignupInfo) (string, error) {
	newUserUuid := utils.GenerateUserUUID()

	// db에 사용자 정보 저장
	// platformInfo와 signupInfo를 사용
	oauthUserInfo, err := s.OauthSessionsRepo.GetOauthInfoBySessionUuid(ctx, oauthSessionUuid)
	if err != nil {
		return "", err
	}
	newUser := &model.User{
		UserUuid:      newUserUuid,
		OAuthPlatform: oauthUserInfo.OauthPlatform,
		OAuthID:       oauthUserInfo.OauthId,
		Email:         signupInfo.Nickname,
		Nickname:      signupInfo.Nickname,
		Bio:           signupInfo.Nickname,
		ProfileImage:  signupInfo.ProfileImgUrl,
	}

	s.UsersRepo.Signup(ctx, newUser)

	return newUserUuid, nil

}

func (s *UserService) IsAvailableNickname(ctx context.Context, nickname string) (bool, error) {
	isAvailableNickname, err := s.UsersRepo.IsAvailableNickname(ctx, nickname)
	if err != nil {
		return false, errors.New("fail to check nickname")
	}

	return isAvailableNickname, nil
}

func (s *UserService) IsUserExists(ctx context.Context, oauthUserInfo *dto.OauthUserInfo) (string, error) {
	userUuid, err := s.UsersRepo.IsUserExists(ctx, oauthUserInfo)
	if err != nil {
		return "", errors.New("fail to check oauth user info")
	}
	logger.Debugf("isUserExists user uuid : %s", userUuid)
	return userUuid, nil
}

func (s *UserService) CreateOauthSession(ctx context.Context, oauthUserInfo *dto.OauthUserInfo) (string, error) {

	sessionID := utils.GenerateRandomSessionID()
	expiryAt := time.Now().Add(time.Duration(config.OAUTH_SESSION_EXPIRY) * time.Minute)
	oauthSession := &model.OAuthSession{
		SessionID:     sessionID,
		OAuthPlatform: oauthUserInfo.OauthPlatform,
		OAuthID:       oauthUserInfo.OauthId,
		ExpiresAt:     expiryAt,
	}
	_, err := s.OauthSessionsRepo.CreateOauthSession(ctx, oauthSession)
	if err != nil {
		return "", errors.New("fail to create oauth session")
	}

	return sessionID, nil
}
