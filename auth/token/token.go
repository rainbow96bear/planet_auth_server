package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/grpc_client"
	"github.com/rainbow96bear/planet_auth_server/logger"
	"github.com/rainbow96bear/planet_auth_server/utils"
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

type TokenProvider struct {
	utils.Provider
}

func CreateAccessToken(userUuid string) (string, error) {

	accessClaims := jwt.MapClaims{
		"userUuid":  userUuid,
		"plateform": "kakao",
		"exp":       time.Now().Add(time.Duration(config.ACCESS_TOKEN_EXPIRY_MINUTE) * time.Minute),
		"iat":       time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(config.JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func CreateRefreshToken() string {
	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)
	refreshToken = fmt.Sprintf("%s", refreshToken)
	return refreshToken
}

func DeleteRefreshToken(token string) error {
	reqToken := &pb.Token{
		Token: token,
	}

	dbClient, err := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	if err != nil {
		return err
	}
	_, err = dbClient.ReqDeleteRefreshToken(reqToken)
	if err != nil {
		return err
	}
	return nil
}

func (t *TokenProvider) IssueAccessToken(c *gin.Context) {
	logger.Infof("start to issue access token")
	defer logger.Infof("end to issue access token")

	refreshToken, err := c.Cookie(config.REFRESH_TOKEN_NAME)
	if err != nil {
		logger.Warnf("no %s cookie found", config.REFRESH_TOKEN_NAME)
		return
	}

	dbClient, err := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	if err != nil {
		logger.Errorf("fail to create dbClient ERR[%s]", err.Error())
		return
	}
	logger.Debugf("request to db server for get user uuid : %+v", refreshToken)

	reqRefreshToken := &pb.Token{
		Token: refreshToken,
	}

	resRefreshTokenInfo, err := dbClient.ReqGetRefreshTokenInfo(reqRefreshToken)
	if err != nil {
		logger.Warnf("failed to get refresh token ERR[%s]", err.Error())
	}
	if resRefreshTokenInfo.GetUserUuid() == "" {
		logger.Warnf("fail to get user uuid from refresh token ERR[%s]", err.Error())
		return
	}
	accessToken, err := CreateAccessToken(resRefreshTokenInfo.GetUserUuid())
	if err != nil {
		logger.Warnf("fail to create access token ERR[%s]", err.Error())
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		900, // 15분
		"/",
		"",   // domain
		true, // secure
		true, // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   900,
	})
}

func (t *TokenProvider) UpdateRefreshTokens(c *gin.Context) {
	logger.Infof("start to refresh tokens")
	defer logger.Infof("end to refresh tokens")

	refreshToken, err := c.Cookie(config.REFRESH_TOKEN_NAME)
	if err != nil {
		logger.Warnf("no refresh_token cookie found")
		return
	}

	dbClient, err := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	if err != nil {
		logger.Errorf("fail to create dbClient ERR[%s]", err.Error())
		return
	}
	logger.Debugf("request to db server for get user id : %+v", refreshToken)

	reqRefreshToken := &pb.Token{
		Token: refreshToken,
	}

	resRefreshTokenInfo, err := dbClient.ReqGetRefreshTokenInfo(reqRefreshToken)
	if err != nil {
		logger.Warnf("failed to oauth sign up Error : %+v", err)
	}

	accessToken, err := CreateAccessToken(resRefreshTokenInfo.UserUuid)
	if err != nil {
		logger.Warnf("fail to create access token ERROR : %s", err.Error())
		return
	}

	newRefreshTokenStr := CreateRefreshToken()

	newRefreshToken := &pb.Token{
		UserUuid: resRefreshTokenInfo.UserUuid,
		Token:    newRefreshTokenStr,
		Expiry:   time.Now().Add(time.Duration(config.ACCESS_TOKEN_EXPIRY_MINUTE) * time.Minute).Unix(),
	}

	logger.Debugf("newRefreshToken value : %+v", newRefreshToken)
	_, err = dbClient.ReqUpdateRefreshToken(newRefreshToken)
	if err != nil {
		logger.Warnf("failed to refresh Token : %+v", err)
		return
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
	c.SetCookie(
		"access_token",
		accessToken,
		900, // 15분
		"/",
		"",   // domain
		true, // secure
		true, // httpOnly
	)
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   900,
	})
}
