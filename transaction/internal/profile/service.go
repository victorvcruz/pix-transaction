package profile

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"transaction/internal/cfg"
	"transaction/internal/errutils"
	proto "transaction/proto/v1"
)

type Service interface {
	FindAccountBalance(ctx context.Context, id string, userId string) (float64, error)
	FindReceiver(ctx context.Context, key string) (*Account, error)
	IsAccountActive(ctx context.Context, id string) (bool, error)
}

type service struct {
	account proto.AccountServiceClient
	user    proto.UserServiceClient
	keys    proto.KeysServiceClient
}

func (s *service) FindAccountBalance(ctx context.Context, id string, userId string) (float64, error) {
	account, err := s.account.FindAccount(ctx, &proto.AccountRequest{AccountId: id, UserId: userId})
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

func (s *service) FindReceiver(ctx context.Context, key string) (*Account, error) {
	account, err := s.account.FindByKey(ctx, &proto.FindByKeyRequest{Key: key})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return nil, errutils.ErrInvalidKey
		}

		return nil, err
	}
	return ProtoToAccount(account), nil
}

func (s *service) IsAccountActive(ctx context.Context, id string) (bool, error) {
	isActive, err := s.account.IsAccountActive(ctx, &proto.AccountRequest{AccountId: id})
	if err != nil {
		return false, err
	}
	return isActive.GetValue(), nil
}

func NewService(config cfg.Config) Service {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(config.ProfileConfig.Host, opts...)
	if err != nil {
		panic(err)
	}

	return &service{
		account: proto.NewAccountServiceClient(conn),
		user:    proto.NewUserServiceClient(conn),
		keys:    proto.NewKeysServiceClient(conn),
	}
}