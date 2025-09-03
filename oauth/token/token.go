package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/grpc_client"
	pb "github.com/rainbow96bear/planet_proto"
)

type accessClaims struct {
	UserId string `json:"userid"`
	Claims jwt.RegisteredClaims
}

type DeleteTokenRequest struct {
	UserID uint64 `json:"user_id"`
	Token  string `json:"token"`
}

func CreateAccessToken(userId uint64) (string, error) {

	accessClaims := jwt.MapClaims{
		"sub":       userId,
		"plateform": "kakao",
		"exp":       time.Now().Add(time.Duration(config.ACCESS_TOKEN_EXPIRY_MINUTE) * time.Minute),
		"iat":       time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(config.JWT_SECRET_KEY)
	if err != nil {
		return "", fmt.Errorf("fail to create accesstoken ERROR[%s]", err.Error())
	}

	return accessTokenString, nil
}

func CreateRefreshToken() (string, error) {
	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)
	refreshToken = fmt.Sprintf("kakao-%s", refreshToken)
	return refreshToken, nil
}

func DeleteRefreshToken(userid uint64, token string) error {

	reqToken := &pb.Token{
		UserId: userid,
		Token:  token,
	}

	dbReq := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	_, err := grpc_client.ReqDeleteRefreshToken(dbReq, reqToken)
	if err != nil {
		return err
	}
	return nil
}
