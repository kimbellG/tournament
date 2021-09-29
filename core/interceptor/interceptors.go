package interceptor

import (
	"context"
	"fmt"

	"github.com/kimbellG/tournament/core/handler/kegrpc"

	"github.com/kimbellG/kerror"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Log(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	entryLog := log.WithFields(log.Fields{
		"action":  info.FullMethod,
		"request": req,
	})
	entryLog.Info("received request")

	resp, err := handler(ctx, req)
	if err != nil {
		kerror.ErrorLog(entryLog, err, fmt.Sprintf("%v failed", info.FullMethod))
	}

	entryLog.WithField("response", resp).Info("request has been processed")

	return resp, err
}

func Error(next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := next(ctx, req, info, handler)
		if err != nil {
			err = kegrpc.Errorf(err, info.FullMethod)
		}

		return resp, err
	}
}
