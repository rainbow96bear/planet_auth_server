package config

var (
	PORT                          string
	LOG_LEVEL                     int16
	ACCESS_TOKEN_EXPIRY_MINUTE    int
	REFRESH_TOKEN_NAME            string
	REFRESH_TOKEN_EXPIRY_DURATION int

	OAUTH_SESSION_EXPIRY int

	KAKAO_CLIENT_SECRET string
	KAKAO_REDIRECT_URI  string
	KAKAO_REST_API_KEY  string

	JWT_SECRET_KEY     string
	PLANET_CLIENT_ADDR string

	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
)
