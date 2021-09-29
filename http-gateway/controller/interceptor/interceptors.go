package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

func Error(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		return DecodeErrorFromGrpc(err)
	}

	return nil
}
