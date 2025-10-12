package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/auth/token"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/dto"
	"github.com/rainbow96bear/planet_auth_server/grpc_client"
	"github.com/rainbow96bear/planet_auth_server/logger"
	"github.com/rainbow96bear/planet_auth_server/utils"
	pb "github.com/rainbow96bear/planet_proto"
)

func Signup(c *gin.Context) {
	logger.Infof("start signup")
	defer logger.Infof("end signup")
	oauthSession, err := c.Cookie("signup_session")
	if err != nil {
		logger.Errorf("fail to get signup_session ERR[%s]", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
		return
	}

	dbClient, err := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	if err != nil {
		logger.Errorf("fail to create dbClient ERR[%s]", err.Error())
		return
	}
	oauthSessionRequest := &pb.GetOauthSessionRequest{
		SessionId: oauthSession,
	}
	resPlatformInfo, err := dbClient.ReqGetPlatformInfoBySession(oauthSessionRequest)
	logger.Debugf("response platform info : %+v", resPlatformInfo)
	// session에 id를 저장

	// 얻은 정보와 받은 body로 db에 넣기
	var req dto.SignupRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Errorf("fail to bind req form ERR[%s]", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}
	logger.Debugf("bind req : %+v", req)

	// 사용자 uuid 만들기
	newUserUuid := utils.GenerateUserUUID()
	file, err := c.FormFile("profile_image")
	var imageURL string
	if err == nil {
		// 파일이 업로드되었으면 서버에 저장하거나 S3 같은 외부 스토리지에 업로드
		savePath := fmt.Sprintf("./profile/image/%s.jpg", newUserUuid)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}

		// 실제 서비스에서는 업로드 후 URL을 생성해서 DB에 넣음
		imageURL = fmt.Sprintf("%s/profile/image/%s.jpg", config.PLANET_CLIENT_ADDR, newUserUuid)
	} else {
		// 파일이 없으면 기본 이미지 사용
		imageURL = fmt.Sprintf("%s/profile/image/default", config.PLANET_CLIENT_ADDR)
	}

	newUserInfo := &pb.UserInfo{
		UserUuid:      newUserUuid,
		OauthPlatform: resPlatformInfo.GetPlatform(),
		OauthId:       resPlatformInfo.GetPlatformId(),
		Email:         req.Email,
		Nickname:      req.Nickname,
		ProfileImage:  imageURL,
		Bio:           req.Bio,
	}

	redirectUrl := fmt.Sprintf("%s/login/callback", config.PLANET_CLIENT_ADDR)

	_, err = dbClient.ReqOauthSignUp(newUserInfo)
	if err != nil {
		logger.Warnf("failed to oauth sign up ERR[%s]", err.Error())
		errorRedirect := fmt.Sprintf("%s?status=error&code=%s", redirectUrl, utils.ERR_DB_REQUEST)
		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"redirect": errorRedirect,
		})
		return
	}

	// 성공
	refreshToken := token.CreateRefreshToken()
	if err != nil {
		logger.Warnf("fail to issue tokens ERR[%s]", err.Error())
	}
	reqRefreshToken := &pb.Token{
		UserUuid: newUserInfo.UserUuid,
		Token:    refreshToken,
		Expiry:   time.Now().Add(3 * 24 * time.Hour).Unix(),
	}

	_, err = dbClient.ReqUpdateRefreshToken(reqRefreshToken)
	if err != nil {
		logger.Warnf("failed to refresh Token ERR[%s]", err.Error())
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

	successRedirect := fmt.Sprintf("%s?status=success", redirectUrl)
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"refresh_token": refreshToken,
		"redirect":      successRedirect,
	})
}

func NicknameAvailable(c *gin.Context) {
	logger.Infof("start nickname available check")
	defer logger.Infof("end nickname available check")

	nickname := c.Query("nickname")
	if len(nickname) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"available": false,
			"error":     "nickname must be at least 2 characters",
		})
		return
	}

	// DB gRPC 클라이언트 생성
	dbClient, err := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	if err != nil {
		logger.Errorf("fail to connect db grpc ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"available": false,
			"error":     "db connection failed",
		})
		return
	}

	// DB 서버에 닉네임 중복 확인 요청
	req := &pb.CheckNicknameRequest{
		Nickname: nickname,
	}

	res, err := dbClient.ReqCheckNicknameAvailable(req)
	if err != nil {
		logger.Errorf("fail to request nickname check ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"available": false,
			"error":     "nickname check failed",
		})
		return
	}

	// DB 서버 응답 전달
	c.JSON(http.StatusOK, gin.H{
		"available": res.Available,
	})
}
