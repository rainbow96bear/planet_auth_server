package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/authInit"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_auth_server/internal/bootstrap"
	"github.com/rainbow96bear/planet_auth_server/internal/router"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
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

	db, err := authInit.InitDB()
	if err != nil {
		logger.Errorf("failed to initialize database: %s", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	handlers := bootstrap.InitHandlers(db)

	r := router.SetupRouter(func(r *gin.Engine) {
		for _, h := range handlers {
			h.RegisterRoutes(r)
		}
	})

	authServerPort := fmt.Sprintf(":%s", config.PORT)

	r.Run(authServerPort)

}
