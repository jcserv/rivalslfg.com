package log

import (
	"github.com/getsentry/sentry-go"
	"github.com/jcserv/rivalslfg/internal/utils/env"
	"go.uber.org/zap/zapcore"
)

func initSentry() error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              env.GetString("SENTRY_DSN", ""),
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Environment:      env.GetString("ENVIRONMENT", "dev"),
	})
}

func sentryHook(entry zapcore.Entry, fields []zapcore.Field) error {
	if entry.Level >= zapcore.ErrorLevel {
		hub := sentry.CurrentHub()
		event := sentry.NewEvent()
		event.Level = sentry.LevelError
		event.Message = entry.Message
		event.Timestamp = entry.Time
		event.Platform = "go"
		event.Logger = entry.LoggerName
		event.Environment = env.GetString("ENVIRONMENT", "dev")
		event.ServerName = "rivalsapi-lfg"

		event.Tags = make(map[string]string)
		event.Extra = make(map[string]interface{})
		event.Contexts = make(map[string]sentry.Context)

		requestData := &sentry.Request{
			Headers: make(map[string]string),
			Env:     make(map[string]string),
		}

		user := sentry.User{}
		for _, field := range fields {
			switch field.Key {
			case "method":
				requestData.Method = field.String
			case "path":
				requestData.URL = field.String
			case "remote_addr":
				requestData.Env["REMOTE_ADDR"] = field.String
			case "user_agent":
				requestData.Headers["User-Agent"] = field.String
			case "request_id":
				event.Tags["request_id"] = field.String
				event.EventID = sentry.EventID(field.String)
			case "transaction_id":
				event.Transaction = field.String
			case "user":
				user.ID = field.String
				event.User = user
			case "os":
				event.Contexts["os"] = sentry.Context{"name": field.String}
			case "device":
				event.Contexts["device"] = sentry.Context{"type": field.String}
			default:
				event.Extra[field.Key] = field.Interface
			}
		}

		if requestData.Method != "" {
			event.Request = requestData
		}

		hub.CaptureEvent(event)
	}
	return nil
}
