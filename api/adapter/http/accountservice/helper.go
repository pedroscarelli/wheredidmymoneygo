package accountservice

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/middleware"
)

var errMissingClaims = errors.New("error getting token claims")

func GetAccountIDFromToken(r *http.Request) (string, error) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*middleware.Claims)
	if !ok || claims.AccountID == "" {
		return "", errMissingClaims
	}

	return claims.AccountID, nil
}

// authorizedAccountID returns the {id} from the path only when it belongs to
// the authenticated account, so users can't read or delete other accounts.
func authorizedAccountID(r *http.Request) (string, int) {
	tokenAccountID, err := GetAccountIDFromToken(r)
	if err != nil {
		return "", http.StatusUnauthorized
	}

	if mux.Vars(r)["id"] != tokenAccountID {
		return "", http.StatusForbidden
	}

	return tokenAccountID, http.StatusOK
}
