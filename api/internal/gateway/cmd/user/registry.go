package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"sync"

	v1 "github.com/and-period/marche/api/internal/gateway/user/v1/handler"
	"github.com/and-period/marche/api/proto/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type registry struct {
	v1 v1.APIV1Handler
}

type options struct {
	logger *zap.Logger
}

type option func(opts *options)

type gRPCClient struct {
	user user.UserServiceClient
}

func withLogger(logger *zap.Logger) option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func newRegistry(conf *config, opts ...option) (*registry, error) {
	// オプション設定の取得
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}

	// gRPC Clientの設定
	grpcCli, err := newGRPCClient(conf)
	if err != nil {
		return nil, err
	}

	// Handlerの設定
	v1Params := &v1.Params{
		Logger:      dopts.logger,
		WaitGroup:   &sync.WaitGroup{},
		UserService: grpcCli.user,
	}

	return &registry{
		v1: v1.NewAPIV1Handler(v1Params),
	}, nil
}

func newGRPCClient(conf *config) (*gRPCClient, error) {
	opts, err := newGRPCOptions(conf)
	if err != nil {
		return nil, err
	}

	if conf.ProxyHost != "" {
		conf.UserServiceURL = conf.ProxyHost
	}

	userConn, err := grpc.Dial(conf.UserServiceURL, opts...)
	if err != nil {
		return nil, err
	}

	return &gRPCClient{
		user: user.NewUserServiceClient(userConn),
	}, nil
}

func newGRPCOptions(conf *config) ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	var cred credentials.TransportCredentials
	if conf.GRPCInsecure {
		cred = insecure.NewCredentials()
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred = credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
	}
	opts = append(opts, grpc.WithTransportCredentials(cred))

	return opts, nil
}
