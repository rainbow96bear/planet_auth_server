package grpc_client

import (
	"context"
	"time"

	"github.com/rainbow96bear/planet_db_server/logger"
	pb "github.com/rainbow96bear/planet_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DBClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

func NewDBClient(addr string) (*DBClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("grpc server did not connect: %v", err)
		return nil, err
	}
	return &DBClient{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (d *DBClient) ReqGetUserInfoByPlatformInfo(platform *pb.PlatformInfo) (*pb.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.client.GetUserInfoByPlatformInfo(ctx, platform)
}

func (d *DBClient) ReqOauthSignUp(userInfo *pb.UserInfo) (*pb.SignUpResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.client.OauthSignUp(ctx, userInfo)
}

func (d *DBClient) ReqUpdateRefreshToken(token *pb.Token) (*pb.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.client.UpdateRefreshToken(ctx, token)
}

func (d *DBClient) ReqGetRefreshTokenInfo(token *pb.Token) (*pb.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.client.GetRefreshTokenInfo(ctx, token)
}

func (d *DBClient) ReqDeleteRefreshToken(token *pb.Token) (*pb.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.client.DeleteRefreshToken(ctx, token)
}

func (d *DBClient) ReqSaveOauthSession(oauthSession *pb.OauthSession) (*pb.SaveOauthSessionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.client.SaveOauthSession(ctx, oauthSession)
}

func (d *DBClient) ReqGetPlatformInfoBySession(oauthSession *pb.GetOauthSessionRequest) (*pb.PlatformInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.client.GetPlatformInfoBySession(ctx, oauthSession)
}

func (d *DBClient) ReqCheckNicknameAvailable(nickname *pb.CheckNicknameRequest) (*pb.CheckNicknameResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.client.CheckNicknameAvailable(ctx, nickname)
}
