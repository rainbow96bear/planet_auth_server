package dto

type SignupRequest struct {
	Nickname string `form:"nickname" binding:"required"`
	Bio      string `form:"bio"`
	Email    string `form:"email"`
}
