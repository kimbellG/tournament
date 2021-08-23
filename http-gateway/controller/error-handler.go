package controller

import (
	"errors"

	"github.com/kimbellG/kerror"
	"github.com/kimbellG/tournament/core/handler/kegrpc/errorpb"
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

	for _, d := range gerr.Details() {
		switch info := d.(type) {
		case errorpb.ErrorHandler:
			return getError(info)
		default:
		}
	}

	return kerror.New(gerr.Err(), UnmarshalStatusCode(gerr.Code()))
}

func getError(err errorpb.ErrorHandler) error {
	switch e := err.Error.(type) {
	case *errorpb.ErrorHandler_Kerror:
		return kerror.New(errors.New(e.Kerror.GetMsg()), kerror.StatusCode(e.Kerror.GetCode()))
	default:
		return kerror.New(errors.New("unexcepted error"), kerror.Unknown)
	}
}

func UnmarshalStatusCode(code codes.Code) kerror.StatusCode {
	kcode, ok := grpcToKe[code]
	if ok {
		return kcode
	}

	return kerror.Unknown
}
