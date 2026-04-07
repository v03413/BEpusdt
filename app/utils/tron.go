package utils

import (
	"context"
	"errors"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

// NewTronGrpcClient 目前很多 TRON API 基本都有QPS限制，过于激进或者过于保守的重试策略都不合适
func NewTronGrpcClient(apiNode string, apiKey []string) (*grpc.ClientConn, error) {
	apiNode = strings.TrimSpace(apiNode)
	if apiNode == "" {
		return nil, errors.New("tron api node address is empty")
	}

	if !strings.Contains(apiNode, "://") {
		apiNode = "passthrough:///" + apiNode
	}

	opts := []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				MaxDelay:   30 * time.Second,
				Multiplier: 1.5,
			},
			MinConnectTimeout: 1 * time.Minute,
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	if len(apiKey) > 0 {
		opts = append(opts,
			grpc.WithUnaryInterceptor(tronGridApiKeyUnaryInterceptor(apiKey)),
			grpc.WithStreamInterceptor(tronGridApiKeyStreamInterceptor(apiKey)),
		)
	}

	return grpc.NewClient(apiNode, opts...)
}

func tronGridApiKeyUnaryInterceptor(apiKeys []string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		i := time.Now().UnixMicro() % int64(len(apiKeys))
		ctx = addTronGridApiKeyToContext(ctx, apiKeys[i])
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func tronGridApiKeyStreamInterceptor(apiKeys []string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		i := time.Now().UnixMicro() % int64(len(apiKeys))
		ctx = addTronGridApiKeyToContext(ctx, apiKeys[i])
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func addTronGridApiKeyToContext(ctx context.Context, apiKey string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "TRON-PRO-API-KEY", apiKey, "tron-pro-api-key", apiKey)
}
