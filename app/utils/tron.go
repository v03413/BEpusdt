package utils

import (
	"context"
	"errors"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func NewTronGrpcClient(apiNode, apiKey string) (*grpc.ClientConn, error) {
	apiNode = strings.TrimSpace(apiNode)
	if apiNode == "" {
		return nil, errors.New("tron api node address is empty")
	}

	if !strings.Contains(apiNode, "://") {
		apiNode = "passthrough:///" + apiNode
	}

	opts := []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.Config{BaseDelay: 1 * time.Second, MaxDelay: 30 * time.Second, Multiplier: 1.5},
			MinConnectTimeout: 1 * time.Minute,
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if apiKey != "" {
		opts = append(opts,
			grpc.WithUnaryInterceptor(tronGridApiKeyUnaryInterceptor(apiKey)),
			grpc.WithStreamInterceptor(tronGridApiKeyStreamInterceptor(apiKey)),
		)
	}

	return grpc.NewClient(apiNode, opts...)
}

func tronGridApiKeyUnaryInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = addTronGridApiKeyToContext(ctx, apiKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func tronGridApiKeyStreamInterceptor(apiKey string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = addTronGridApiKeyToContext(ctx, apiKey)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func addTronGridApiKeyToContext(ctx context.Context, apiKey string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "TRON-PRO-API-KEY", apiKey, "tron-pro-api-key", apiKey)
}
