package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger   *zerolog.Logger
	appName  = "app"
	env      = os.Getenv("APP_ENV")
	initOnce sync.Once
	mu       sync.Mutex
)

// initLogger initializes the logger only once
func initLogger() {
	if env == "development" {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			FormatLevel: func(l interface{}) string {
				return fmt.Sprintf("[%s]", l)
			},
		}
		l := zerolog.New(consoleWriter).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger = &l
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		fileLogger := &lumberjack.Logger{
			Filename:   fmt.Sprintf("logs/%s.log", appName),
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
		}

		multi := io.MultiWriter(os.Stderr, fileLogger)
		l := zerolog.New(multi).With().Timestamp().Logger()
		logger = &l
	}
}

// getLogger ensures initialization, no lock needed for reading
func getLogger() *zerolog.Logger {
	initOnce.Do(initLogger)
	return logger
}

// SetAppName rebuilds logger with a new name (locks to avoid races)
func SetAppName(name string) {
	mu.Lock()
	defer mu.Unlock()
	appName = name

	// rebuild logger for new name
	initLogger()
}

func SetLogLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

// Expose logger methods
func Debug() *zerolog.Event { return getLogger().Debug() }
func Info() *zerolog.Event  { return getLogger().Info() }
func Warn() *zerolog.Event  { return getLogger().Warn() }
func Error() *zerolog.Event { return getLogger().Error() }
func Fatal() *zerolog.Event { return getLogger().Fatal() }
func Panic() *zerolog.Event { return getLogger().Panic() }
