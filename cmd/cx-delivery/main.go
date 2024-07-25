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
	Log.Info().Msgf("delivery agent started")

	// // subscribe to topics where channel=processor
	// jsonPayload := `{"channel_in": ["processor"], "channel_not_in": [""]}`
	// data, _, err := util.PostUrl(
	// 	viper.GetString("APP_BACKEND_URL")+"/api/v1/superuser/channeltopics",
	// 	[]byte(jsonPayload),
	// 	util.SetConfigAuth,
	// 	nil)
	// if err != nil {
	// 	Log.Error().Msgf("Error fetching processor destination")
	// }

	// var topicsToSubscribe []ChannelTopic
	// if err = json.Unmarshal([]byte(data), &topicsToSubscribe); err != nil {
	// 	Log.Error().Msgf("Error Unmarshalling topicsToSubscribe by delivery agent")
	// }

	// // init topicMan, but topicMan needs processorFunc() :(

	// // get write keys for channel=postgres
	// jsonPayload = `{"channel_in": ["postgres"], "channel_not_in": [""]}`
	// data, _, err = util.PostUrl(
	// 	viper.GetString("APP_BACKEND_URL")+"/api/v1/superuser/channeltopics",
	// 	[]byte(jsonPayload),
	// 	util.SetConfigAuth,
	// 	nil)
	// if err != nil {
	// 	Log.Error().Msgf("Error fetching processor destination")
	// }

	// var channelTopics []ChannelTopic
	// if unmarshalErr := json.Unmarshal([]byte(data), &channelTopics); unmarshalErr != nil {
	// 	Log.Error().Msgf("Error Unmarshalling event: %v", unmarshalErr)
	// 	return
	// }

	// for _, ct := range channelTopics {
	// 	headerItems := map[string]string{"Content-Type": "application/json", "X-Write-Key": ct.WriteKey}
	// 	Log.Info().Msgf("************************* making 3049 with %v", headerItems)
	// 	eventBytes, _ := json.Marshal(event)
	// 	_, _, err = util.PostUrl("http://localhost:3049/api/s/s2s/event", eventBytes, nil, headerItems)
	// 	if err != nil {
	// 		Log.Error().Msgf("error sending request to Jitsu: %v", err)
	// 	}
	// }
	// initialize environment & config variables
	env.InitConfig()
	Log.Info().Msg(viper.GetString("APP_BACKEND_URL"))
	Log.Info().Msg(viper.GetString("KAFKA_BROKER"))

	// initialize ConfigStore
	jsonPayload := `{"channel_in": ["processor"], "channel_not_in": ["x", "y"]}`
	cs, err := configstore.Init(jsonPayload)
	if err != nil {
		Log.Fatal().Msg(err.Error())
	}

	// initialize Broker
	topicMan, err := InitBroker(delivery)
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
