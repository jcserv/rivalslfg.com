package log

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jcserv/rivalslfg/internal/utils/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

func Init(isProd bool) *zap.Logger {
	if isProd {
		if err := initSentry(); err != nil {
			panic(err)
		}
		defer sentry.Flush(2 * time.Second)
	}

	config := zap.NewProductionConfig()
	if !isProd {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.TimeKey = "time"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	l, _ := config.Build()

	options := []zap.Option{zap.AddCallerSkip(1)}
	if isProd {
		sentryCore := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.RegisterHooks(core, func(entry zapcore.Entry) error {
				fields := []zapcore.Field{}
				return sentryHook(entry, fields)
			})
		})
		options = append(options, sentryCore)
	}

	logger = l.WithOptions(options...)
	return logger
}

func WithRequest(l *zap.Logger, r *http.Request) *zap.Logger {
	userAgent := r.UserAgent()
	os, device := parseUserAgent(userAgent)

	return l.With(
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", userAgent),
		zap.String("request_id", r.Header.Get("X-Request-ID")),
		zap.String("user", r.Header.Get("X-User-ID")),
		zap.String("device", device),
		zap.String("os", os),
		zap.String("url", r.URL.String()),
	)
}

func parseUserAgent(userAgent string) (os, device string) {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "windows"):
		os = "windows"
	case strings.Contains(ua, "mac os"):
		os = "macos"
	case strings.Contains(ua, "linux"):
		os = "linux"
	case strings.Contains(ua, "android"):
		os = "android"
	case strings.Contains(ua, "ios"):
		os = "ios"
	default:
		os = "unknown"
	}

	switch {
	case strings.Contains(ua, "mobile"):
		device = "mobile"
	case strings.Contains(ua, "tablet"):
		device = "tablet"
	default:
		device = "desktop"
	}
	return
}

func GetLogger(ctx context.Context) *zap.Logger {
	isProd := env.GetString("ENVIRONMENT", "dev") == "production"

	if logger == nil {
		return Init(isProd)
	}
	return logger
}

func Debug(ctx context.Context, msg string) {
	if !env.GetBool("DEBUG", false) {
		return
	}
	GetLogger(ctx).Debug(msg)
}

func Info(ctx context.Context, msg string) {
	GetLogger(ctx).Info(msg)
}

func Warn(ctx context.Context, msg string) {
	GetLogger(ctx).Warn(msg)
}

func Error(ctx context.Context, msg string) {
	GetLogger(ctx).Error(msg)
}

func Fatal(ctx context.Context, msg string) {
	GetLogger(ctx).Fatal(msg)
}
