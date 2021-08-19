package controller

import (
	"github.com/kimbellG/kerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var grpcToKe = map[codes.Code]kerror.StatusCode{
	codes.InvalidArgument:    kerror.BadRequest,
	codes.NotFound:           kerror.NotFound,
	codes.FailedPrecondition: kerror.SQLConstraintError,
	codes.Internal:           kerror.SQLQueryError,
	codes.Aborted:            kerror.SQLTransactionError,
	codes.Unknown:            kerror.Unknown,
}

func decodeGrpcError(err error) error {
	gerr, ok := status.FromError(err)
	if !ok {
		return kerror.New(err, kerror.Unknown)
	}

	return kerror.New(gerr.Err(), UnmarshalStatusCode(gerr.Code()))
}

func UnmarshalStatusCode(code codes.Code) kerror.StatusCode {
	kcode, ok := grpcToKe[code]
	if ok {
		return kcode
	}

	return kerror.Unknown
}
