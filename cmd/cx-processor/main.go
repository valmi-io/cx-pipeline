/*
 * Copyright (c) 2024 valmi.io <https://github.com/valmi-io>
 *
 * Created Date: Wednesday, July 17th 2024, 6:11:58 pm
 * Author: Rajashekar Varkala @ valmi.io
 */

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/valmi-io/cx-pipeline/internal/configstore"
	"github.com/valmi-io/cx-pipeline/internal/env"
	. "github.com/valmi-io/cx-pipeline/internal/log"
	. "github.com/valmi-io/cx-pipeline/internal/msgbroker"
)

func main() {
	// initialize environment & config variables
	env.InitConfig()
	Log.Info().Msg(viper.GetString("APP_BACKEND_URL"))
	Log.Info().Msg(viper.GetString("KAFKA_BROKER"))

	// initialize ConfigStore
	jsonPayload := `{"channel_in": ["chatbox"], "channel_not_in": ["x", "y"]}`
	cs, err := configstore.Init(jsonPayload)
	if err != nil {
		Log.Fatal().Msg(err.Error())
	}

	// initialize Broker
	topicMan, err := InitBroker(processor)
	if err != nil {
		Log.Fatal().Msg(err.Error())
	}

	// Connect ConfigStore and Broker
	cs.AttachTopicMan(topicMan)

	cleanupChan := make(chan bool)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanupChan <- true
	}()

	<-cleanupChan
	Log.Info().Msg("Received an interrupt signal, shutting down...")
	cs.Close()       // close ConfigStore to stop refreshing IFTTTs
	topicMan.Close() // close broker topic management

}
