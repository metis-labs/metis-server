package types

import "context"

// UserIDFromCtx returns the userID from the given context.
func UserIDFromCtx(ctx context.Context) string {
	return ctx.Value(userIDKey).(string)
}

// CtxWithUserID creates a new context with the given userID.
func CtxWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
