package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/auth/token"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/logger"
	"github.com/rainbow96bear/planet_auth_server/oauth/kakao"
	"github.com/rainbow96bear/planet_auth_server/router"
	"google.golang.org/grpc"
)

// go build -ldflags "-X main.Mode=prod -X main.Version=1.0.0 -X main.GitCommit=$(git rev-parse HEAD)" -o user_service_prod .

var (
	Mode      string
	Version   string
	GitCommit string
)

func init() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version: %s\nCommit: %s\n", Version, GitCommit)
		os.Exit(0)
	}
	Mode = "dev"
	fmt.Printf("user_service Start \nVersion : %s \nGit Commit : %s\n", Version, GitCommit)
	fmt.Printf("Build Mode : %s\n", Mode)
	config.InitConfig(Mode)
	logger.SetLevel(config.LOG_LEVEL)
}

func main() {

	kakaoClient := &oauth.KakaoClient{
		RestApiKey:   config.KAKAO_REST_API_KEY,
		RedirectUrl:  config.KAKAO_REDIRECT_URI,
		ClientSecret: config.KAKAO_CLIENT_SECRET,
    }

    // TokenService 초기화 (JWT 발급용)
    tokenService := &service.TokenService{
		AccessTokenExpiry : config.ACCESS_TOKEN_EXPIRY_MINUTE
		JwtSecretKey : config.JWT_SECRET_KEY
	}

    // KakaoHandler 생성
    kakaoHandler := &handler.KakaoHandler{
        KakaoClient:  kakaoClient,
        TokenService: tokenService,
    }

	r := router.SetupRouter(
		func(r *gin.Engine) { routes.RegisterKakaoOauthRoutes(r, kakaoHandler) },
	)

	authServerPort := fmt.Sprintf(":%s", config.PORT)

	r.Run(authServerPort)

}
