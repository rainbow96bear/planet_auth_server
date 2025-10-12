package service

import (
	"github.com/rainbow96bear/planet_auth_server/utils"
	pb "github.com/rainbow96bear/planet_proto"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	providers map[string]utils.Provider
}

func NewUserService(providers map[string]utils.Provider) *UserService {
	return &UserService{providers: providers}
}
