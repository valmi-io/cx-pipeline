package main

import (
	"encoding/json"

	"github.com/spf13/viper"
	. "github.com/valmi-io/cx-pipeline/internal/log"
	util "github.com/valmi-io/cx-pipeline/internal/util"
)

type ChannelTopic struct {
	LinkID     string `json:"link_id"`
	WriteKey   string `json:"write_key"`
	storefront string
	channel    string
}

func delivery(msg string) {
	Log.Info().Msgf("processing msg %v", msg)

	var event map[string]interface{}
	if unmarshalErr := json.Unmarshal([]byte(msg), &event); unmarshalErr != nil {
		Log.Error().Msgf("Error Unmarshalling event: %v", unmarshalErr)
		return
	}

	jsonPayload := `{"channel_in": ["postgres"], "channel_not_in": [""]}`
	data, _, err := util.PostUrl(
		viper.GetString("APP_BACKEND_URL")+"/api/v1/superuser/channeltopics",
		[]byte(jsonPayload),
		util.SetConfigAuth,
		nil)
	if err != nil {
		Log.Error().Msgf("Error fetching processor destination")
	}

	var channelTopics []ChannelTopic
	if unmarshalErr := json.Unmarshal([]byte(data), &channelTopics); unmarshalErr != nil {
		Log.Error().Msgf("Error Unmarshalling event: %v", unmarshalErr)
		return
	}

	for _, ct := range channelTopics {
		headerItems := map[string]string{"Content-Type": "application/json", "X-Write-Key": ct.WriteKey}
		eventBytes, _ := json.Marshal(event)
		_, _, err = util.PostUrl("http://localhost:3049/api/s/s2s/event", eventBytes, nil, headerItems)
		if err != nil {
			Log.Error().Msgf("error sending request to Jitsu: %v", err)
		}
	}

}
