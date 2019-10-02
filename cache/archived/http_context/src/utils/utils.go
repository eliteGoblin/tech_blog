package utils

import (
	"context"
	"net/http"
)

type key int

const requestIDKey key = 0

func NewContextWithRequestID(ctx context.Context, req *http.Request) context.Context {
	reqID := req.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = "frank-random-1234"
	}

	return context.WithValue(ctx, requestIDKey, reqID)
}

func RequestIDFromContext(ctx context.Context) string {
	if ctx != nil && ctx.Value(requestIDKey) != nil {
		return ctx.Value(requestIDKey).(string)
	}
	return "nil context"
}
