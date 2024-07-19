package configstore

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	. "github.com/valmi-io/cx-pipeline/internal/log"
	. "github.com/valmi-io/cx-pipeline/internal/msgbroker"
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
	// Pass to matchChannelState
	// if success, switch to NewChannelTopics
}
type ChannelTopic struct {
	LinkID string    `json:"link_id"`
	WriteKey   string `json:"write_key"`
	storeID string
	channel string
}

type ChannelTopics struct {
	Channels []ChannelTopic;
	done chan bool
}

func (t *ChannelTopic) constructTopic() string{
	return fmt.Sprintf("in.%v", t.LinkID)
}

func matchChannelState(newCT *ChannelTopics, currentCT *ChannelTopics) bool {
	if newCT == nil || newCT.Channels == nil {
		Log.Info().Msg("New Channel Topics is either empty or null")
	} else{
		Log.Info().Msg("Subscribing to new Channels")
		for _, nt  := range newCT.Channels {
			if currentCT != nil && currentCT.Channels != nil {
				for _,ct := range currentCT.Channels {
					if ct.LinkID == nt.LinkID {
						break // already subscribed
					}
				}
			}
			SubscribeTopic(nt.constructTopic()) // subscribe
		}
	}
	
	if currentCT == nil || currentCT.Channels == nil {
		Log.Info().Msg("Current Channel Topics is either empty or null")
	} else{
		Log.Info().Msg("Unsubscribing from old Channels")
        for _,ct := range currentCT.Channels {
            if newCT!= nil && newCT.Channels!= nil {
                for _,nt := range newCT.Channels {
                    if ct.LinkID == nt.LinkID {
                        break // already unsubscribed
                    }
                }
            }
            UnsubscribeTopic(ct.constructTopic()) // unsubscribe
        }
	}
	return true
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

