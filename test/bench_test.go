package test

import (
	"errors"
	"github.com/Calyr3x/zeroslog"
	"io"
	"log/slog"
	"testing"
	"time"

	_ "github.com/rs/zerolog/log"

	"github.com/sirupsen/logrus"
)

var (
	fullMethod = "/conveyor.NotificationService/GetAmountUnreadNotifications"
	dur        = 320 * time.Microsecond
	version    = 21
	someFloat  = 123.123
	errTest    = errors.New("this is an error")
)

func BenchmarkSlog_Info(b *testing.B) {
	slogLogger := slog.New(
		slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	benchLog(b, func() {
		slogLogger.Info("gRPC call succeeded",
			"method", fullMethod,
			"duration", dur,
		)
	})
}

func BenchmarkZeroSLog_Info(b *testing.B) {
	logger := slog.New(zeroslog.New(
		zeroslog.WithTimeFormat("2006-01-02 15:04:05.000 -07:00"),
		zeroslog.WithOutput(io.Discard),
		zeroslog.WithColors(),
		zeroslog.WithMinLevel(0),
	))

	benchLog(b, func() {
		logger.Info("gRPC call succeeded",
			"method", fullMethod,
			"duration", dur,
		)
	})
}

func BenchmarkLogrus_Info(b *testing.B) {
	l := logrus.New()
	l.SetOutput(io.Discard)
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000 -07:00",
	})

	benchLog(b, func() {
		l.WithFields(logrus.Fields{
			"method":   fullMethod,
			"duration": dur,
		}).Info("gRPC call succeeded")
	})
}

func BenchmarkSlog_Error(b *testing.B) {
	slogLogger := slog.New(
		slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	benchLog(b, func() {
		slogLogger.Error("gRPC call succeeded",
			"method", fullMethod,
			"duration", dur,
			"err", errTest,
		)
	})
}

func BenchmarkZeroSLog_Error(b *testing.B) {
	logger := slog.New(zeroslog.New(
		zeroslog.WithTimeFormat("2006-01-02 15:04:05.000 -07:00"),
		zeroslog.WithOutput(io.Discard),
		zeroslog.WithColors(),
	))

	benchLog(b, func() {
		logger.Error("gRPC call succeeded",
			"method", fullMethod,
			"duration", dur,
			"err", errTest,
		)
	})
}

func BenchmarkLogrus_Error(b *testing.B) {
	l := logrus.New()
	l.SetOutput(io.Discard)
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000 -07:00",
	})

	benchLog(b, func() {
		l.WithFields(logrus.Fields{
			"method":   fullMethod,
			"duration": dur,
		}).Error("err", errTest)
	})
}

func benchLog(b *testing.B, fn func()) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn()
	}
}
