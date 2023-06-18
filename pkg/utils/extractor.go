package utils

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5/request"
)

// CookieExtractor extracts a token from request cookies
type CookieExtractor []string

func (e CookieExtractor) ExtractToken(req *http.Request) (string, error) {
	for _, arg := range e {
		if cookie, _ := req.Cookie(arg); cookie != nil {
			return cookie.Value, nil
		}
	}

	return "", request.ErrNoTokenInRequest
}
