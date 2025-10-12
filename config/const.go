package config

var (
	PORT                          string
	GRPC_PORT                     string
	LOG_LEVEL                     int16
	DB_SERVER_ADDR                string
	DB_GRPC_SERVER_ADDR           string
	ACCESS_TOKEN_EXPIRY_MINUTE    int16
	REFRESH_TOKEN_NAME            string
	REFRESH_TOKEN_EXPIRY_DURATION int

	KAKAO_CLIENT_SECRET string
	KAKAO_REDIRECT_URI  string
	KAKAO_REST_API_KEY  string

	JWT_SECRET_KEY     string
	PLANET_CLIENT_ADDR string
)
