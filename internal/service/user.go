package service

type UserService struct {
	// AccessTokenExpiry 
	// JwtSecretKey string
	ProfileImgSavePath string
	UserRepo *repository.UserRepository
}
func (s *UserService)Signup(oauthSessionUuid string, signupInfo SignupInfo)(error){
	newUserUuid := utils.GenerateUserUUID()

	// db에 사용자 정보 저장
	// platformInfo와 signupInfo를 사용
	platformInfo, err := s.SessionRepo.GetOauthSessionInfo(oauthSessionUuid)

	
}

func (s *UserService)IsAvailableNickname(nickname string) (bool, error) {
	//db에 nickname 중복 검사

	if s.UserRepo.IsNicknameTaken(user.Nickname) {
        return false, errors.New("nickname already taken")
    }

	return true, nil
}

func (s *UserService)IsUserExists(platformInfo PlatformInfo) bool {

}