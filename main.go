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

	// caCert, err := os.ReadFile("ca.crt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// caCertPool := x509.NewCertPool()
	// if !caCertPool.AppendCertsFromPEM(caCert) {
	// 	logger.Errorf("failed to append CA cert")
	// }

	// 서버 인증서/개인키 로드 (self-signed or CA signed)
	// serverCert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	// if err != nil {
	// 	logger.Errorf("failed to load server key pair: %v", err)
	// }

	kakaoOauthProvider := &kakao.OauthProvider{
		RestApiKey:   config.KAKAO_REST_API_KEY,
		RedirectUrl:  config.KAKAO_REDIRECT_URI,
		ClientSecret: config.KAKAO_CLIENT_SECRET,
	}

	authTokenProvider := &token.TokenProvider{}
	r := router.SetupRouter(
		func(r *gin.Engine) { router.RegisterKakaoOauthRoutes(r, kakaoOauthProvider) },
		func(r *gin.Engine) { router.RegisterTokenRoutes(r, authTokenProvider) },
		func(r *gin.Engine) { router.RegisterSignupRoutes(r) },
		// router.RegisterPostRoutes,
	)
	authServerPort := fmt.Sprintf(":%s", config.PORT)
	// go func() {
	// 	logger.Infof("Starting HTTPS server on %s", authServerPort)
	// 	if err := r.RunTLS(authServerPort, "server.crt", "server.key"); err != nil {
	// 		logger.Errorf("failed to start HTTPS server: %v", err)
	// 	}
	// }()

	r.Run(authServerPort)
	authServerGrpcPort := fmt.Sprintf(":%s", config.GRPC_PORT)
	lis, err := net.Listen("tcp", authServerGrpcPort)
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
	}

	// creds := credentials.NewTLS(&tls.Config{
	// 	Certificates: []tls.Certificate{serverCert},
	// 	ClientCAs:    caCertPool,
	// 	// 개발 환경에서는 클라이언트 인증 강제하지 않음
	// 	ClientAuth: tls.NoClientCert,
	// })

	// grpcServer := grpc.NewServer(grpc.Creds(creds))
	grpcServer := grpc.NewServer()

	// userService := service.NewUserService(map[string]utils.Provider{
	// 	"kakao": utils.Provider(kakaoOauthProvider),
	// })
	// pb.RegisterUserServiceServer(grpcServer, userService)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Errorf("failed to serve: %v", err)
	}

}
