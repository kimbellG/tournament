package kegrpc

import (
	"errors"
	"fmt"

	"github.com/kimbellG/tournament/core/handler/kegrpc/errorpb"

	"github.com/kimbellG/kerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var keToGrpcDict = map[kerror.StatusCode]codes.Code{
	kerror.InvalidID:                   codes.InvalidArgument,
	kerror.BadRequest:                  codes.InvalidArgument,
	kerror.NotFound:                    codes.NotFound,
	kerror.TournamentDoesntExists:      codes.NotFound,
	kerror.UserDoesntExists:            codes.NotFound,
	kerror.SQLConstraintError:          codes.FailedPrecondition,
	kerror.SQLQueryError:               codes.Internal,
	kerror.SQLPrepareStatementError:    codes.Internal,
	kerror.SQLScanError:                codes.Internal,
	kerror.SQLExecutionError:           codes.Internal,
	kerror.SQLTransactionError:         codes.Aborted,
	kerror.SQLTransactionBeginError:    codes.Aborted,
	kerror.SQLTransactionRoolbackError: codes.Aborted,
	kerror.SQLTransactionCommitError:   codes.Aborted,
	kerror.Unknown:                     codes.Unknown,
}

func Newf(code kerror.StatusCode, format string, args ...interface{}) error {
	grpcCode := MarshalStatusCode(code)
	st := status.New(grpcCode, grpcCode.String())

	ds, err := st.WithDetails(
		&errorpb.ErrorHandler{
			Error: &errorpb.ErrorHandler_Kerror{
				Kerror: &errorpb.Kerror{
					Code: int64(code),
					Msg:  fmt.Sprintf(format, args...),
				},
			},
		},
	)
	if err != nil {
		return st.Err()
	}

	return ds.Err()

}

func Errorf(err error, format string, args ...interface{}) error {
	args = append(args, err)

	if trnt := (kerror.Error{}); errors.As(err, &trnt) {
		return Newf(trnt.StatusCode(), format+": %v", args...)
	}

	return Newf(kerror.Unknown, format+": %v", args...)
}

func MarshalStatusCode(code kerror.StatusCode) codes.Code {
	grpcCode, ok := keToGrpcDict[code]
	if !ok {
		return codes.Unknown
	}

	return grpcCode
}
