package ctxutil

import "context"

const keyValidatedForm key = "validatedForm"

type key string

func ValidatedForm[T any](ctx context.Context) (T, bool) {
	form, ok := ctx.Value(keyValidatedForm).(T)
	return form, ok
}

func SetValidatedForm[T any](ctx context.Context, form T) context.Context {
	return context.WithValue(ctx, keyValidatedForm, form)
}
