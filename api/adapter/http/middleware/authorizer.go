package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/response"
	domain "github.com/pedroscarelli/wheredidmymoneygo/core/domain"
)

type Claims struct {
	AccountID string `json:"sub"`
	jwt.RegisteredClaims
}

type contextKey string

var UserClaimsKey contextKey = "userClaims"

func JwtAuthorizer(jwtSecret []byte, accountUseCase domain.AccountUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || tokenStr == "" {
				respondUnauthorized(w)
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
				return jwtSecret, nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())

			if err != nil || !token.Valid || claims.AccountID == "" {
				respondUnauthorized(w)
				return
			}

			if _, err := accountUseCase.GetByID(claims.AccountID); err != nil {
				if errors.Is(err, domain.ErrAccountNotFound) {
					respondUnauthorized(w)
					return
				}
				response.InternalError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func respondUnauthorized(w http.ResponseWriter) {
	response.Error(w, http.StatusUnauthorized, "unauthorized")
}
