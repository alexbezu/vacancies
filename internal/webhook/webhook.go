package webhook

import (
	"context"

	"github.com/sirupsen/logrus"
)

// LogWebhook is a webhook that logs messages to the console.
type LogWebhook struct {
	log *logrus.Logger
}

// NewLogWebhook creates a new log webhook.
func NewLogWebhook(log *logrus.Logger) *LogWebhook {
	return &LogWebhook{log: log}
}

// Send sends the message to the console.
func (w *LogWebhook) Send(ctx context.Context, message string) error {
	w.log.Info(message)
	return nil
}
