package handler

import (
	"errors"
	"net/http"

	"github.com/kimbellG/kerror"
)

var httpCodeDict = map[kerror.StatusCode]int{
	kerror.BadRequest:          http.StatusBadRequest,
	kerror.InternalServerError: http.StatusBadRequest,
	kerror.InvalidID:           http.StatusBadRequest,
	kerror.NotFound:            http.StatusNotFound,

	kerror.SQLConstraintError: http.StatusBadRequest,

	kerror.SQLQueryError:            http.StatusInternalServerError,
	kerror.SQLExecutionError:        http.StatusInternalServerError,
	kerror.SQLPrepareStatementError: http.StatusInternalServerError,
	kerror.SQLScanError:             http.StatusInternalServerError,

	kerror.SQLTransactionError:         http.StatusInternalServerError,
	kerror.SQLTransactionBeginError:    http.StatusInternalServerError,
	kerror.SQLTransactionCommitError:   http.StatusInternalServerError,
	kerror.SQLTransactionRoolbackError: http.StatusInternalServerError,

	kerror.TournamentDoesntExists: http.StatusNotFound,
	kerror.UserDoesntExists:       http.StatusNotFound,
}

func decodeStatusCode(err error) int {
	if terr := (kerror.Error{}); errors.As(err, &terr) {
		if httpCode, ok := httpCodeDict[terr.StatusCode()]; ok {
			return httpCode
		}
	}

	return http.StatusInternalServerError
}
