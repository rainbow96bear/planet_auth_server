package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/google/uuid"
)

type Provider interface {
}

// generateRandomSessionID는 안전한 랜덤 32바이트 세션 ID를 생성합니다.
func GenerateRandomSessionID() string {
	bytes := make([]byte, 32) // 32바이트 = 256비트
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("failed to generate random session ID: %v", err)
	}
	return hex.EncodeToString(bytes) // 64자리 hex 문자열
}

func GenerateUserUUID() string {
	newUUID := uuid.New() // Version 4 UUID (랜덤 기반)
	return newUUID.String()
}
