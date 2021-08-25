package controller

import (
	"github.com/kimbellG/kerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var grpcToErrorCode = map[codes.Code]kerror.StatusCode{
	codes.InvalidArgument:    kerror.BadRequest,
	codes.NotFound:           kerror.NotFound,
	codes.FailedPrecondition: kerror.SQLConstraintError,
	codes.Internal:           kerror.SQLQueryError,
	codes.Aborted:            kerror.SQLTransactionError,
	codes.Unknown:            kerror.Unknown,
}

func decodeErrorFromGrpc(err error) error {
	gerr, ok := status.FromError(err)
	if !ok {
		return kerror.New(err, kerror.Unknown)
	}

	for _, d := range gerr.Details() {
		switch info := d.(type) {
		case errorpb.ErrorHandler:
			return errorFromProto(info)
		default:
		}
	}

	return kerror.New(gerr.Err(), UnmarshalStatusCode(gerr.Code()))
}

func errorFromProto(err errorpb.ErrorHandler) error {
	switch e := err.Error.(type) {
	case *errorpb.ErrorHandler_Kerror:
		return kerror.New(errors.New(e.Kerror.GetMsg()), kerror.StatusCode(e.Kerror.GetCode()))
	default:
		return kerror.New(errors.New("unexcepted error"), kerror.Unknown)
	}
}

func UnmarshalStatusCode(code codes.Code) kerror.StatusCode {
	kcode, ok := grpcToErrorCode[code]
	if ok {
		return kcode
	}

	return kerror.Unknown
}
