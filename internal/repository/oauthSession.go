package repository

import (
	"context"
	"database/sql"
	"fmt"
	"planet_utils/model"
	"planet_utils/pkg/logger"

	"github.com/rainbow96bear/planet_auth_server/dto"
)

type OauthSessionRepository struct {
	DB *sql.DB
}

func (r *OauthSessionRepository) GetOauthInfoBySessionUuid(ctx context.Context, oauthSessionUuid string) (*dto.OauthUserInfo, error) {
	logger.Infof("start to get oauth session info by session uuid")
	defer logger.Infof("end to get oauth session info by session uuid")

	query := `
		SELECT oauth_platform, oauth_id
		FROM oauth_sessions
		WHERE session_id = ?;
	`

	var oauthUserInfo dto.OauthUserInfo

	err := r.DB.QueryRowContext(ctx, query, oauthSessionUuid).Scan(
		&oauthUserInfo.OauthPlatform,
		&oauthUserInfo.OauthId,
	)
	if err != nil {
		logger.Errorf("failed to get oauth session info ERR[%s]", err.Error())
		return nil, err
	}

	logger.Debugf("successfully fetched oauth session info: %+v", oauthUserInfo)
	return &oauthUserInfo, nil
}

func (r *OauthSessionRepository) CreateOauthSession(ctx context.Context, oauthSession *model.OAuthSession) (*model.OAuthSession, error) {
	logger.Infof("start to create oauth session: %+v", oauthSession)
	defer logger.Infof("end to create oauth session")

	query := `
		INSERT INTO oauth_sessions (session_id, oauth_platform, oauth_id, expires_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		oauthSession.SessionID,
		oauthSession.OAuthPlatform,
		oauthSession.OAuthID,
		oauthSession.ExpiresAt,
	)

	if err != nil {
		logger.Errorf("failed to create oauth session ERR[%s]", err.Error())
		return nil, fmt.Errorf("failed to insert oauth session: %w", err)
	}

	logger.Debugf("successfully created oauth session: %+v", oauthSession)
	return oauthSession, nil
}
