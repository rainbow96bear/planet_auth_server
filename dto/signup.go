package dto

type SignupInfo struct {
	Nickname      string `form:"nickname" binding:"required"`
	Bio           string `form:"bio"`
	Email         string `form:"email"`
	ProfileImgUrl string `from:"profileImg"`
}

type OauthUserInfo struct {
	OauthPlatform string
	OauthId       string
}
