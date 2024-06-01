package transport

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func LogPrint(logger *log.Logger, message string) {
	if logger != nil {
		logger.Println(message)
	}
}

func LoggingInterceptorServer(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if logger == nil {
			return handler(ctx, req)
		}

		start := time.Now()

		resp, err := handler(ctx, req)

		logger.Printf("RPC response: %s, Duration: %v, Error: %v", info.FullMethod, time.Since(start), err)

		return resp, err
	}
}

func LoggingInterceptorClient(logger *log.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if logger == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		logger.Printf("RPC request: %s, Duration: %v, Error: %v", method, time.Since(start), err)

		return err
	}
}
