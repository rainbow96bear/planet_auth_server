package planetInit

import (
	"database/sql"
	"fmt"
	"planet_utils/pkg/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rainbow96bear/planet_auth_server/config"
)

var DB *sql.DB

// InitDB initializes and returns a database connection.
func InitDB() (*sql.DB, error) {
	// 환경 변수에서 DB 접속 정보 읽기
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		config.DB_USER,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME,
	)
	logger.Debugf(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// 커넥션 풀 설정 (운영 환경에서 중요)
	db.SetMaxOpenConns(50)           // 최대 연결 개수
	db.SetMaxIdleConns(10)           // 유휴 연결 개수
	db.SetConnMaxLifetime(time.Hour) // 커넥션 재사용 시간

	// 연결 확인
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	logger.Infof("✅ Successfully connected to database [%s:%s/%s]", config.DB_HOST, config.DB_PORT, config.DB_NAME)

	DB = db
	return db, nil
}
