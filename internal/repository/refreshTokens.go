package repository

import (
	"context"
	"database/sql"
	"planet_utils/model"
	"planet_utils/pkg/logger"
)

type RefreshTokensRepository struct {
	DB *sql.DB
}

func (r *RefreshTokensRepository) GetUserUuidByRefreshToken(ctx context.Context, refresh_token string) (string, error) {
	logger.Infof("start to get user uuid by refresh token")
	defer logger.Infof("end to get user uuid by refresh token")

	query := `
		SELECT user_uuid
		FROM refresh_tokens
		WHERE token = ?
	`

	var userUUID string
	err := r.DB.QueryRowContext(ctx, query, refresh_token).Scan(&userUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warnf("no user found for refresh token: %s", refresh_token)
			return "", nil
		}
		logger.Errorf("failed to query user info by refresh token: %v", err)
		return "", err
	}

	logger.Debugf("found user uuid: %s", userUUID)
	return userUUID, nil
}

func (r *RefreshTokensRepository) UpdateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) (*model.RefreshToken, error) {
	logger.Infof("start to update refresh token")
	defer logger.Infof("end to update refresh token")

	query := `
		INSERT INTO refresh_tokens (user_uuid, token, expiry)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			token = VALUES(token),
			expiry = VALUES(expiry)
	`

	_, err := r.DB.ExecContext(ctx, query, refreshToken.UserUUID, refreshToken.Token, refreshToken.Expiry)
	if err != nil {
		logger.Errorf("failed to update refresh token: %v", err)
		return nil, err
	}

	logger.Debugf("refresh token updated for refreshToken: %+v", refreshToken)
	return refreshToken, nil
}

func (r *RefreshTokensRepository) DeleteRefreshToken(ctx context.Context, refreshTokenStr string) error {
	logger.Infof("start to delete refresh token")
	defer logger.Infof("end to delete refresh token")

	queryDelete := `DELETE FROM refresh_tokens WHERE token = ?`
	_, err := r.DB.ExecContext(ctx, queryDelete, refreshTokenStr)
	if err != nil {
		logger.Errorf("failed to delete refresh token: %v", err)
		return err
	}

	logger.Infof("refresh token deleted")
	return nil
}
