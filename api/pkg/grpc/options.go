package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxzap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type options struct {
	logger *zap.Logger
}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewGRPCOptions(opts ...Option) []grpc.ServerOption {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}

	streamInterceptors := grpcStreamServerInterceptors(dopts.logger)
	unaryInterceptors := grpcUnaryServerInterceptors(dopts.logger)

	gopts := []grpc.ServerOption{
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	}
	return gopts
}

/*
 * ServerOptions - StremServerInterceptor
 */
func grpcStreamServerInterceptors(logger *zap.Logger) []grpc.StreamServerInterceptor {
	opts := []grpc_zap.Option{
		grpc_zap.WithDecider(shouldLog),
	}

	interceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_zap.StreamServerInterceptor(logger, opts...),
		grpc_recovery.StreamServerInterceptor(),
	}

	return interceptors
}

/*
 * ServerOptions - UnaryServerInterceptor
 */
func grpcUnaryServerInterceptors(logger *zap.Logger) []grpc.UnaryServerInterceptor {
	opts := []grpc_zap.Option{
		grpc_zap.WithDecider(shouldLog),
	}

	interceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(logger, opts...),
		accessLogUnaryServerInterceptor(),
		grpc_recovery.UnaryServerInterceptor(),
	}

	return interceptors
}

func accessLogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		res, err := handler(ctx, req)

		// ヘルスチェック時はログ出力のスキップ
		if !shouldLog(info.FullMethod, err) {
			return res, err
		}

		clientIP := ""
		if p, ok := peer.FromContext(ctx); ok {
			clientIP = p.Addr.String()
		}

		requestID := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if id, ok := md["x-request-id"]; ok {
				requestID = strings.Join(id, ",")
			}
		}

		userAgent := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if u, ok := md["user-agent"]; ok {
				userAgent = strings.Join(u, ",")
			}
		}

		// Request / Response Messages
		reqParams := map[string]interface{}{}
		if p, ok := req.(proto.Message); ok {
			reqParams, _ = filterParams(p)
		}

		resParams := map[string]interface{}{}
		if p, ok := res.(proto.Message); ok {
			resParams, _ = filterParams(p)
		}

		ds := getErrorDetails(err)

		grpc_ctxzap.AddFields(
			ctx,
			zap.String("request.client_ip", clientIP),
			zap.String("request.request_id", requestID),
			zap.String("request.user_agent", userAgent),
			zap.Reflect("request.content", reqParams),
			zap.Reflect("response.content", resParams),
			zap.Reflect("error.details", ds),
		)

		return res, err
	}
}

func shouldLog(fullMethodName string, err error) bool {
	return err != nil || fullMethodName != "/grpc.health.v1.Health/Check"
}

func filterParams(pb proto.Message) (map[string]interface{}, error) {
	fields := []string{"auth", "password"}

	bs, err := protojson.Marshal(pb)
	if err != nil {
		return nil, fmt.Errorf("jsonpb serializer failed: %w", err)
	}

	bj := make(map[string]interface{})
	_ = json.Unmarshal(bs, &bj) // ignore error here.

	var toFilter []string
	for k := range bj {
		for i := range fields {
			if !strings.Contains(strings.ToLower(k), fields[i]) {
				continue
			}
			toFilter = append(toFilter, k)
		}
	}

	for _, k := range toFilter {
		bj[k] = "<FILTERED>"
	}

	return bj, nil
}

func getErrorDetails(err error) interface{} {
	if err == nil {
		return ""
	}

	s, ok := status.FromError(err)
	if !ok {
		return ""
	}

	// TODO: 配列に1つしか値入れないようにしてるからいったんこれで
	for _, detail := range s.Details() {
		switch v := detail.(type) {
		case *errdetails.LocalizedMessage:
			return v.GetMessage()
		case *errdetails.BadRequest:
			return v.GetFieldViolations()
		}
	}

	return ""
}
