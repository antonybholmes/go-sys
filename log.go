package sys

import (
	"io"

	"os"

	"github.com/rs/zerolog"
)

func NewFileLog(f string) (*zerolog.Logger, error) {
	logFile, err := os.OpenFile(f, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}

	defer logFile.Close()

	log := zerolog.New(io.MultiWriter(os.Stdout, logFile)).With().Timestamp().Logger()

	return &log, nil
}
