package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"module-go/internal/cfg"
	"os"
	"path/filepath"
	"strconv"
)

func Init() {
	// Enable pretty logging if debug mode is enabled
	if cfg.Get().Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	// Enable unix time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Enable caller
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = log.With().Caller().Logger()

	// Set log level
	logLevel := cfg.Get().LogLevel

	var zeroLogLevel zerolog.Level
	if err := zeroLogLevel.UnmarshalText([]byte(logLevel)); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal log level")
	}

	zerolog.SetGlobalLevel(zeroLogLevel)

	log.Info().Msg("Logger successfully initialized")
}