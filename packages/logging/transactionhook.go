package logging

import (
	"dietician.local/packages/constants"
	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type TransactionHook struct{}

func (t *TransactionHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (t *TransactionHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		if requestID, ok := entry.Context.Value(constants.RequestIDKey).(string); ok {
			entry.Data["requestId"] = requestID
		}
		if userID, ok := entry.Context.Value(constants.UserIDKey).(string); ok {
			entry.Data["userId"] = userID
		}
		if clientID, ok := entry.Context.Value(constants.ClientIDKey).(string); ok {
			entry.Data["clientId"] = clientID
		}
		if requestFields, ok := entry.Context.Value(constants.RequestLogFieldsKey).(logrus.Fields); ok {
			entry.Data["request"] = requestFields
		}
		if txn := newrelic.FromContext(entry.Context); txn != nil {
			logcontext.AddLinkingMetadata(entry.Data, txn.GetLinkingMetadata())
		}
	}

	return nil
}

func NewTransactionHook() *TransactionHook {
	return &TransactionHook{}
}
