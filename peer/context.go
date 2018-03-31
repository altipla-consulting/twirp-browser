package peer

import (
	"net/http"

	"golang.org/x/net/context"
)

type key int

const requestKey key = 0

func RequestWithContext(r *http.Request) context.Context {
	return context.WithValue(r.Context(), requestKey, r)
}

func RequestFromContext(ctx context.Context) *http.Request {
	return ctx.Value(requestKey).(*http.Request)
}

func AuthorizationFromContext(ctx context.Context) string {
	r := RequestFromContext(ctx)
	return r.Header.Get("Authorization")
}
