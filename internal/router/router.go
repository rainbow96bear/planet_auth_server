package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter는 gin.Engine을 초기화하고
// 각 route 등록 함수를 받아 router를 완성해줍니다.
func SetupRouter(routeFuncs ...func(*gin.Engine)) *gin.Engine {
	// Gin 엔진 초기화
	r := gin.New()

	// 기본 미들웨어 추가
	// r.Use(gin.Logger())   // 요청 로그 출력
	r.Use(gin.Recovery()) // panic 발생 시 복구

	// CORS 설정
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 필요시 도메인 제한 가능
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 전달받은 route 등록 함수 실행
	for _, fn := range routeFuncs {
		fn(r)
	}

	// health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return r
}
