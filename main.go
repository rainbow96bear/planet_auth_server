package main

import (
	"flag"
	"fmt"
	"os"
	"planet_utils/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/external/oauthClient"
	"github.com/rainbow96bear/planet_auth_server/internal/handler"
	"github.com/rainbow96bear/planet_auth_server/internal/repository"
	"github.com/rainbow96bear/planet_auth_server/internal/router"
	"github.com/rainbow96bear/planet_auth_server/internal/routes"
	"github.com/rainbow96bear/planet_auth_server/internal/service"
	"github.com/rainbow96bear/planet_auth_server/planetInit"
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

	db, err := planetInit.InitDB()
	if err != nil {
		logger.Errorf("failed to initialize database: %s", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	kakaoClient := &oauthClient.KakaoClient{
		RestApiKey:   config.KAKAO_REST_API_KEY,
		RedirectUrl:  config.KAKAO_REDIRECT_URI,
		ClientSecret: config.KAKAO_CLIENT_SECRET,
	}

	userRepo := &repository.UserRepository{
		DB: db,
	}

	oauthRepo := &repository.OauthSessionRepository{
		DB: db,
	}

	refreshTokensRepo := &repository.RefreshTokensRepository{
		DB: db,
	}
	// UserService 초기화
	userService := &service.UserService{
		ProfileImgSavePath: "/profile/image",
		UserRepo:           userRepo,
		OauthSessionRepo:   oauthRepo,
	}

	// TokenService 초기화 (JWT 발급용)
	tokenService := &service.TokenService{
		AccessTokenExpiry:  config.ACCESS_TOKEN_EXPIRY_MINUTE,
		RefreshTokenName:   config.REFRESH_TOKEN_NAME,
		RefreshTokenExpiry: config.REFRESH_TOKEN_EXPIRY_DURATION,
		JwtSecretKey:       config.JWT_SECRET_KEY,

		RefreshTokensRepo: refreshTokensRepo,
	}

	// KakaoHandler 생성
	kakaoHandler := &handler.KakaoHandler{
		KakaoClient:  kakaoClient,
		UserService:  userService,
		TokenService: tokenService,
		Platform:     "kakao",
	}

	tokenHandler := &handler.TokenHandler{}

	userHandler := &handler.UserHandler{
		UserService:  userService,
		TokenService: tokenService,
	}

	r := router.SetupRouter(
		func(r *gin.Engine) { routes.RegisterKakaoOauthRoutes(r, kakaoHandler) },
		func(r *gin.Engine) { routes.RegisterTokenRoutes(r, tokenHandler) },
		func(r *gin.Engine) { routes.RegisterUserRoutes(r, userHandler) },
	)

	authServerPort := fmt.Sprintf(":%s", config.PORT)

	r.Run(authServerPort)

}
