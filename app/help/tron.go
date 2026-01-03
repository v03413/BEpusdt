package help

import (
	"context"
	"time"

	"github.com/v03413/bepusdt/app/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// NewTronGrpcClient 创建带有API key认证的TRON gRPC客户端连接
// 如果配置了tron_api_key，则会自动添加到gRPC metadata中
func NewTronGrpcClient() (*grpc.ClientConn, error) {
	var grpcParams = grpc.ConnectParams{
		Backoff:           backoff.Config{BaseDelay: 1 * time.Second, MaxDelay: 30 * time.Second, Multiplier: 1.5},
		MinConnectTimeout: 1 * time.Minute,
	}

	// 创建拦截器来添加API key
	apiKey := conf.GetTronApiKey()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithConnectParams(grpcParams))
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 如果配置了API key，添加unary和stream拦截器
	if apiKey != "" {
		// Unary拦截器（用于普通RPC调用）
		unaryInterceptor := func(
			ctx context.Context,
			method string,
			req, reply interface{},
			cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker,
			opts ...grpc.CallOption,
		) error {
			// 添加API key到metadata，同时支持多种格式
			// TronGrid 支持这些 header 名称
			ctx = metadata.AppendToOutgoingContext(ctx,
				"TRON-PRO-API-KEY", apiKey,
				"tron-pro-api-key", apiKey, // 小写版本
			)
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		// Stream拦截器（用于流式RPC调用）
		streamInterceptor := func(
			ctx context.Context,
			desc *grpc.StreamDesc,
			cc *grpc.ClientConn,
			method string,
			streamer grpc.Streamer,
			opts ...grpc.CallOption,
		) (grpc.ClientStream, error) {
			// 添加API key到metadata，同时支持多种格式
			ctx = metadata.AppendToOutgoingContext(ctx,
				"TRON-PRO-API-KEY", apiKey,
				"tron-pro-api-key", apiKey, // 小写版本
			)
			return streamer(ctx, desc, cc, method, opts...)
		}

		opts = append(opts, grpc.WithUnaryInterceptor(unaryInterceptor))
		opts = append(opts, grpc.WithStreamInterceptor(streamInterceptor))
	}

	return grpc.NewClient(conf.GetTronGrpcNode(), opts...)
}
