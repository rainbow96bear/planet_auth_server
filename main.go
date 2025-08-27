package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/logger"
	"github.com/rainbow96bear/planet_auth_server/oauth/kakao"
	"github.com/rainbow96bear/planet_auth_server/router"
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

	kakaoOauth := &kakao.Provider{
		RestApiKey:   config.KAKAO_REST_API_KEY,
		RedirectUrl:  config.KAKAO_REDIRECT_URI,
		ClientSecret: config.KAKAO_CLIENT_SECRET,
	}

	r := router.SetupRouter(
		func(r *gin.Engine) { router.RegisterKakaoOauthRoutes(r, kakaoOauth) },
		// router.RegisterPostRoutes,
	)
	authServerPort := fmt.Sprintf(":%s", config.PORT)
	r.Run(authServerPort)

}
