package handler

import (
	"fmt"
	"net/http"
	"planet_utils/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/dto"
	"github.com/rainbow96bear/planet_auth_server/internal/service"
	"github.com/rainbow96bear/planet_auth_server/utils"
)

type UserHandler struct {
	UserService  *service.UserService
	TokenService *service.TokenService
}

func (h *UserHandler) Signup(c *gin.Context) {
	logger.Infof("start to signup")
	defer logger.Infof("end to signup")
	ctx := c.Request.Context()
	oauthSession, err := c.Cookie("signup_session")
	if err != nil {
		logger.Errorf("fail to get signup_session ERR[%s]", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
		return
	}

	var req *dto.SignupInfo
	if err := c.ShouldBind(&req); err != nil {
		logger.Errorf("fail to bind req form ERR[%s]", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}
	logger.Debugf("bind req : %+v", req)

	file, err := c.FormFile("profile_image")
	var imageURL string
	if err == nil {
		// 파일이 업로드되었으면 서버에 저장하거나 S3 같은 외부 스토리지에 업로드
		savePath := fmt.Sprintf("./%s/%s.jpg", h.UserService.ProfileImgSavePath, req.Nickname)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}

		// 실제 서비스에서는 업로드 후 URL을 생성해서 DB에 넣음
		imageURL = fmt.Sprintf("%s/%s/%s.jpg", config.PLANET_CLIENT_ADDR, h.UserService.ProfileImgSavePath, req.Nickname)
	} else {
		// 파일이 없으면 기본 이미지 사용
		imageURL = fmt.Sprintf("%s/%s/default", config.PLANET_CLIENT_ADDR, h.UserService.ProfileImgSavePath)
	}

	// TODO : imageURL 처리
	req.ProfileImgUrl = imageURL
	userUuid, err := h.UserService.Signup(ctx, oauthSession, req)

	redirectUrl := fmt.Sprintf("%s/login/callback", config.PLANET_CLIENT_ADDR)
	if err != nil {
		logger.Warnf("failed to oauth sign up ERR[%s]", err.Error())
		errorRedirect := fmt.Sprintf("%s?status=error&code=%s", redirectUrl, utils.ERR_DB_REQUEST)
		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"redirect": errorRedirect,
		})
		return
	}

	refreshToken, err := h.TokenService.IssueRefreshToken(ctx, userUuid)
	if err != nil {
		logger.Warnf("failed to refresh Token ERR[%s]", err.Error())
		errorRedirect := fmt.Sprintf("%s?status=error&code=%s", redirectUrl, utils.ERR_DB_REQUEST)
		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"redirect": errorRedirect,
		})
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

	successRedirect := fmt.Sprintf("%s?status=success", redirectUrl)
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"refresh_token": refreshToken,
		"redirect":      successRedirect,
	})
}

func (h *UserHandler) NicknameCheck(c *gin.Context) {
	logger.Infof("start to nickname check")
	defer logger.Infof("end to nickname check")
	ctx := c.Request.Context()
	nickname := c.Query("nickname")
	if len(nickname) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"available": false,
			"error":     "nickname must be at least 2 characters",
		})
		return
	}

	isAvailableNickname, err := h.UserService.IsAvailableNickname(ctx, nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"available": false,
			"error":     "nickname already taken",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available": isAvailableNickname,
	})
}
