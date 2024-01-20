package ctx

import "context"

const keyRequestID key = "requestID"

type key string

func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(keyRequestID).(string)

	return requestID
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, keyRequestID, requestID)
}
