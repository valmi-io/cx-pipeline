package configstore

import (
	"time"

	"github.com/spf13/viper"
	. "github.com/valmi-io/cx-pipeline/internal/log"
	util "github.com/valmi-io/cx-pipeline/internal/util"
)

func fetchChannelTopics() (string, error) {
	data, respCode, err := util.PostUrl(
		viper.GetString("APP_BACKEND_URL") + "/api/v1/superuser/channeltopics", 
		[]byte("{\"channel_in\": [\"x\",\"y\"], \"channel_not_in\": [\"x\",\"y\"]}"),
		util.SetConfigAuth)

	Log.Info().Msg(data)
	Log.Info().Msgf("%v", respCode)
	if err!= nil {
		Log.Info().Msg(err.Error())
    }
	return data, err

	// PARSE the response and store it in the ChannelTopics struct 
}

type ChannelTopics struct {
	Channels []struct {
        linkID int    `json:"link_id"`
        writeKey   string `json:"write_key"`
		storeID string
		channel string
    }
	done chan bool
}


func initChannelTopics() (*ChannelTopics, error) {
	d,_:= time.ParseDuration(viper.GetString("CONFIG_REFRESH_INTERVAL"))
	ticker := time.NewTicker(d)
	channelTopics := ChannelTopics{done: make(chan bool)}
	go func() {
		for {
            select {
            case <- channelTopics.done:
                return
            case t := <-ticker.C:
				Log.Debug().Msgf("ChannelTopics Refresh Tick at %v", t)
				fetchChannelTopics()
            }
        }
	}()
	return &channelTopics, nil
}

func (ct *ChannelTopics) Close(){
	ct.done <- true
}

