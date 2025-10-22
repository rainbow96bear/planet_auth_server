package repository

import (
	"context"
	"database/sql"
	"fmt"
	"planet_utils/model"
	"planet_utils/pkg/logger"

	"github.com/rainbow96bear/planet_auth_server/dto"
)

type UsersRepository struct {
	DB *sql.DB
}

func (r *UsersRepository) Signup(ctx context.Context, userInfo *model.User) (string, error) {
	logger.Infof("start signup process for nickname: %s", userInfo.Nickname)
	defer logger.Infof("end signup process for nickname: %s", userInfo.Nickname)

	query := `
		INSERT INTO users (user_uuid, oauth_platform, oauth_id, email, nickname, profile_image, bio)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.DB.ExecContext(ctx, query,
		userInfo.UserUuid,
		userInfo.OAuthPlatform,
		userInfo.OAuthID,
		userInfo.Email,
		userInfo.Nickname,
		userInfo.ProfileImage,
		userInfo.Bio,
	)

	if err != nil {
		logger.Errorf("failed to insert user ERR[%s]", err.Error())
		return "", fmt.Errorf("failed to insert user: %w", err)
	}

	logger.Infof("successfully inserted user: %s", userInfo.UserUuid)
	return userInfo.UserUuid, nil
}

func (r *UsersRepository) IsAvailableNickname(ctx context.Context, nickname string) (bool, error) {
	logger.Infof("start to checking if nickname is available: %s", nickname)
	defer logger.Infof("end to checking nickname: %s", nickname)

	query := `SELECT COUNT(*) FROM users WHERE nickname = ?`
	var count int

	err := r.DB.QueryRowContext(ctx, query, nickname).Scan(&count)
	if err != nil {
		logger.Errorf("failed to check nickname availability ERR[%s]", err.Error())
		return false, err
	}

	available := count == 0
	if available {
		logger.Debugf("nickname '%s' is available", nickname)
	} else {
		logger.Debugf("nickname '%s' is already taken", nickname)
	}

	return available, nil
}

func (r *UsersRepository) IsUserExists(ctx context.Context, oauthUserInfo *dto.OauthUserInfo) (string, error) {
	logger.Infof("start to get user uuid")
	defer logger.Infof("end to get user uuid")

	query := `
        SELECT user_uuid
        FROM users
        WHERE oauth_platform = ? and oauth_id = ?;
    `

	var userUuid string
	err := r.DB.QueryRowContext(ctx, query, oauthUserInfo.OauthPlatform, oauthUserInfo.OauthId).Scan(&userUuid)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debugf("user not found: %+v", oauthUserInfo)
			return "", nil // 조회 결과 없으면 사용 가능
		}
		logger.Errorf("failed to get user uuid ERR[%s]", err.Error())
		return "", err
	}

	logger.Debugf("successfully got user uuid: %s", userUuid)
	return userUuid, nil // 조회 결과 있으면 이미 존재
}
