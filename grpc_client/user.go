package grpc_client

import (
	"context"
	"time"

	"github.com/rainbow96bear/planet_db_server/logger"
	pb "github.com/rainbow96bear/planet_proto"
	"google.golang.org/grpc"
)

func NewDBClient(addr string) pb.UserServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Errorf("grpc server did not connect: %v", err)
	}
	return pb.NewUserServiceClient(conn)
}

func ReqOauthSignUp(client pb.UserServiceClient, userInfo *pb.UserInfo) (*pb.SignUpResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if res, err := client.OauthSignUp(ctx, userInfo); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
