package config

import (
	"time"

	"github.com/rs/zerolog"
)

func Init() error {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return nil
}
