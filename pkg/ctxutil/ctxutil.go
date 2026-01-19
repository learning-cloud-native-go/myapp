package ctxutil

import "context"

const (
	keyRequestID     key = "requestID"
	keyValidatedForm key = "validatedForm"
)

type key string

func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(keyRequestID).(string)

	return requestID
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, keyRequestID, requestID)
}

func ValidatedForm[T any](ctx context.Context) (T, bool) {
	form, ok := ctx.Value(keyValidatedForm).(T)
	return form, ok
}

func SetValidatedForm[T any](ctx context.Context, form T) context.Context {
	return context.WithValue(ctx, keyValidatedForm, form)
}
