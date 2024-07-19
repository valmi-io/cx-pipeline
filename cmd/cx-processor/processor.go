package main

import (
	. "github.com/valmi-io/cx-pipeline/internal/log"
)

func processor(msg string) {
	Log.Info().Msgf("processing msg %v", msg)
}
