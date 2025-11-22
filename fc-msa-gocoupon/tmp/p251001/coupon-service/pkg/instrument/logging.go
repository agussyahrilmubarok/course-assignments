package instrument

import "context"

func GetRequestID(ctx context.Context) string {
	if rid, ok := ctx.Value(RequestIDKey).(string); ok {
		return rid
	}
	return ""
}

func GetTraceID(ctx context.Context) string {
	if tid, ok := ctx.Value(TraceIDKey).(string); ok {
		return tid
	}
	return ""
}
