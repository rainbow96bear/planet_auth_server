package repository

import (
    "database/sql"
    "your_module_name/internal/model"
)

type OauthSessionRepository struct {
    DB *sql.DB
}

func (r *OauthSessionRepository) GetOauthSessionInfo(oauthSessionUuid string) (*model.PlatformInfo, error) {
    _, err := r.DB.Exec("INSERT INTO users (id, nickname, email) VALUES (?, ?, ?)", user.ID, user.Nickname, user.Email)
    return err
}