package accountservice

import (
	"context"

	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/middleware"
)

func contextWithAccount(ctx context.Context, accountID string) context.Context {
	return context.WithValue(ctx, middleware.UserClaimsKey, &middleware.Claims{AccountID: accountID})
}
