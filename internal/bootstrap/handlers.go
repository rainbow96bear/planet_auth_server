package bootstrap

import (
	"database/sql"

	"github.com/rainbow96bear/planet_auth_server/internal/handler"
	"github.com/rainbow96bear/planet_auth_server/internal/repository"
	"github.com/rainbow96bear/planet_auth_server/internal/service"
	"github.com/rainbow96bear/planet_utils/pkg/router"
)

func InitHandlers(db *sql.DB) map[string]router.RouteRegistrar {
	userRepo := &repository.UsersRepository{DB: db}
	refreshTokensRepo := &repository.RefreshTokensRepository{DB: db}
	oauthSessionRepo := &repository.OauthSessionsRepository{DB: db}

	userService := &service.UserService{UsersRepo: userRepo, OauthSessionsRepo: oauthSessionRepo}
	tokenService := &service.TokenService{RefreshTokensRepo: refreshTokensRepo}

	return map[string]router.RouteRegistrar{
		"kakao": handler.NewKakaoHandler(userService, tokenService),
		"token": handler.NewTokenHandler(tokenService),
		"user":  handler.NewUserHandler(userService, tokenService),
	}
}
