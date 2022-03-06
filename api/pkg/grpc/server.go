package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server interface {
	Serve() error
	Stop()
}

type gRPCServer struct {
	server *grpc.Server
	lister net.Listener
}

// NewGRPCServer - gRPCサーバーの生成
func NewGRPCServer(server *grpc.Server, port int64) (Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("grpc: failed to listen port: %w", err)
	}
	s := &gRPCServer{server: server, lister: lis}
	s.registerHealthServer()
	return s, nil
}

// Serve - サーバーの起動
func (s *gRPCServer) Serve() error {
	return s.server.Serve(s.lister)
}

// Stop - サーバーの停止
func (s *gRPCServer) Stop() {
	s.server.GracefulStop()
}

// registerHealthServer - ヘルスチェックエンドポイントの追加
func (s *gRPCServer) registerHealthServer() {
	RegisterHealthServer(s.server)
}
