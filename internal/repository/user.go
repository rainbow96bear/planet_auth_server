package repository

import (
    "database/sql"
    "your_module_name/internal/model"
)

type UserRepository struct {
    DB *sql.DB
}

func (r *UserRepository) GetOauthSessionInfo(oauthSessionUuid string) (*model.PlatformInfo, error) {
    _, err := r.DB.Exec("INSERT INTO users (id, nickname, email) VALUES (?, ?, ?)", user.ID, user.Nickname, user.Email)
    return err
}

func (r *UserRepository) IsAvailableNickname(nickname string) bool {
    var count int
    r.DB.QueryRow("SELECT COUNT(*) FROM users WHERE nickname = ?", nickname).Scan(&count)
    return count > 0
}