package logger

import (
	"github.com/sirupsen/logrus"
)

func WithRequestID(requestID string) *logrus.Entry {
	return Log.WithField("request_id", requestID)
}
