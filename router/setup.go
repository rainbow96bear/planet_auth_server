package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/logger"
)

func SetupRouter(registerFns ...func(*gin.Engine)) *gin.Engine {
	logger.Infof("start set up router")
	defer logger.Infof("end set up router")

	logger.Debugf("register funcs %+v", registerFns)
	// r := gin.Default()
	r := gin.New()

	// 전달받은 라우트 등록 함수 실행
	for _, register := range registerFns {
		register(r)
	}

	return r
}
