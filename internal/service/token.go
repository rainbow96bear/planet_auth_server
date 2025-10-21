package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"planet_utils/model"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rainbow96bear/planet_auth_server/internal/repository"
)

type TokenService struct {
	AccessTokenExpiry  int
	RefreshTokenName   string
	RefreshTokenExpiry int
	JwtSecretKey       string

	RefreshTokensRepo *repository.RefreshTokensRepository
}

func (s *TokenService) IssueAccessToken(ctx context.Context, refreshToken string) (string, error) {
	userUuid, err := s.RefreshTokensRepo.GetUserUuidByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	accessClaims := jwt.MapClaims{
		"userUuid":  userUuid,
		"plateform": "kakao",
		"exp":       time.Now().Add(time.Duration(s.AccessTokenExpiry) * time.Minute),
		"iat":       time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func (s *TokenService) IssueRefreshToken(ctx context.Context, userUuid string) (string, error) {
	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)
	refreshToken = fmt.Sprintf("%s", refreshToken)

	newRefreshToken := &model.RefreshToken{
		UserUUID: userUuid,
		Token:    refreshToken,
		Expiry:   uint64(s.RefreshTokenExpiry),
	}

	updatedToken, err := s.RefreshTokensRepo.UpdateRefreshToken(ctx, newRefreshToken)
	if err != nil {
		return "", err
	}

	return updatedToken.Token, nil
}

func (s *TokenService) ReissueRefreshToken(ctx context.Context, refreshTokenStr string) (string, error) {
	userUuid, err := s.RefreshTokensRepo.GetUserUuidByRefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return "", err
	}

	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)
	refreshToken = fmt.Sprintf("%s", refreshToken)

	newRefreshToken := &model.RefreshToken{
		UserUUID: userUuid,
		Token:    refreshToken,
		Expiry:   uint64(s.RefreshTokenExpiry),
	}

	updatedToken, err := s.RefreshTokensRepo.UpdateRefreshToken(ctx, newRefreshToken)
	if err != nil {
		return "", err
	}

	return updatedToken.Token, nil
}

func (s *TokenService) RevokeRefreshToken(ctx context.Context, refreshTokenStr string) error {
	err := s.RefreshTokensRepo.DeleteRefreshToken(ctx, refreshTokenStr)
	return err
}
