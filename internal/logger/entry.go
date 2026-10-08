package logger

import (
	"cleaning/internal/common/requestcontext"
	"context"

	"github.com/sirupsen/logrus"
)

func WithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	if requestID := requestcontext.RequestID(ctx); requestID != "" {
		fields["request_id"] = requestID
	}
	return Log.WithFields(fields)
}
