package service

type UserService struct {
	// AccessTokenExpiry 
	// JwtSecretKey string
	ProfileImgSavePath string
}

func (s *UserService)CreateUserInfo(PlatformInfo, SignupInfo) (UserInfo, error){
	newUserUuid := utils.GenerateUserUUID()


}

func (s *UserService)CreateProfileImgFileName(file) (string, error){

	return profileImgName, nil
}