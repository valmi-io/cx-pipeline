package main

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/valmi-io/cx-pipeline/internal/env"
	"github.com/valmi-io/cx-pipeline/internal/log"
)
func main() {
    // initialize config
    env.InitConfig() 
    fmt.Println(viper.Get("APP_BACKEND"))
    fmt.Println(viper.Get("KAFKA_BROKER"))

    // initialize logging configuration 
    log := log.GetLogger()
    log.Debug().Msg("This message appears only when log level set to Debug")
    log.Info().Msg("This message appears when log level set to Debug or Info")

}