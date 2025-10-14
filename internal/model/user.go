package model

import "time"

// User : 회원 정보 구조체
type User struct {
    ID        string    `json:"id"`        // UUID 등
    Nickname  string    `json:"nickname"`  // 닉네임
    Email     string    `json:"email"`     // 이메일
    Password  string    `json:"password"`  // 해시된 비밀번호
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}