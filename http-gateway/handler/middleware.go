package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type AuthenticationMiddleware struct {
	TokenPassword string
	NotAuthPaths  []string
}

func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if amw.isPathWithAuthentication(r.URL.Path) {
			if err := amw.logIn(r); err != nil {
				http.Error(w, "Failed of user authentication: "+err.Error(), http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (amw *AuthenticationMiddleware) isPathWithAuthentication(path string) bool {
	for _, nap := range amw.NotAuthPaths {
		if nap == path {
			return false
		}
	}

	return true
}

func (amw *AuthenticationMiddleware) logIn(r *http.Request) error {
	tkString, err := getAuthTokenString(r)
	if err != nil {
		return fmt.Errorf("get token from header: %v", err)
	}

	tk, err := jwt.ParseWithClaims(tkString, &LogClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(amw.TokenPassword), nil
	})
	if err := ValidateToken(tk, err); err != nil {
		return fmt.Errorf("validate token from header: %v", err)
	}

	return nil
}

func getAuthTokenString(r *http.Request) (string, error) {
	headerValue := r.Header.Get("Authorization")
	if headerValue == "" {
		return "", fmt.Errorf("authorization header doesn't exists")

	}

	tkString, err := parseAuthHeader(headerValue)
	if err != nil {
		return "", fmt.Errorf("invalid format of auth token: %v", err)
	}

	return tkString, nil
}

func parseAuthHeader(value string) (string, error) {
	splitted := strings.Split(value, " ")
	if len(splitted) != 2 {
		return "", fmt.Errorf("header value should have 2 words")
	}

	if splitted[0] != "Bearer" {
		return "", fmt.Errorf("first word of header value should be \"Bearer\"")
	}

	return splitted[1], nil
}

func ValidateToken(tk *jwt.Token, err error) error {
	if !tk.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return errors.New("token expired")
			} else {
				return errors.New("Couldn't handle this token: " + err.Error())
			}
		}

		return errors.New("Couldn't handle this token: " + err.Error())
	}

	return nil
}
