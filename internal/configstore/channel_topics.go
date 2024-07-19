package configstore

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	. "github.com/valmi-io/cx-pipeline/internal/log"
	. "github.com/valmi-io/cx-pipeline/internal/msgbroker"
	util "github.com/valmi-io/cx-pipeline/internal/util"
)

func fetchChannelTopics() (string, error) {
	data, respCode, err := util.PostUrl(
		viper.GetString("APP_BACKEND_URL")+"/api/v1/superuser/channeltopics",
		[]byte("{\"channel_in\": [\"x\",\"y\"], \"channel_not_in\": [\"x\",\"y\"]}"),
		util.SetConfigAuth)

	Log.Info().Msg(data)
	Log.Info().Msgf("%v", respCode)
	if err != nil {
		Log.Info().Msg(err.Error())
	}
	return data, err

	// PARSE the response and store it in the ChannelTopics struct
	// Pass to matchChannelState
	// if success, switch to NewChannelTopics
}

type ChannelTopic struct {
	LinkID   string `json:"link_id"`
	WriteKey string `json:"write_key"`
	storeID  string
	channel  string
}

type ChannelTopics struct {
	Channels []ChannelTopic
	done     chan bool
	topicMan *TopicMan
}

func (t *ChannelTopic) constructTopic() string {
	return fmt.Sprintf("in.%v", t.LinkID)
}

func matchChannelState(newCT *ChannelTopics, currentCT *ChannelTopics) bool {
	if currentCT.topicMan == nil {
		return false
	}

	if newCT == nil || newCT.Channels == nil {
		Log.Info().Msg("New Channel Topics is either empty or null")
	} else {
		Log.Info().Msg("Subscribing to new Channels")
		for _, nt := range newCT.Channels {
			if currentCT != nil && currentCT.Channels != nil {
				for _, ct := range currentCT.Channels {
					if ct.LinkID == nt.LinkID {
						break // already subscribed
					}
				}
			}
			currentCT.topicMan.SubscribeTopic(nt.constructTopic()) // subscribe
		}
	}

	if currentCT == nil || currentCT.Channels == nil {
		Log.Info().Msg("Current Channel Topics is either empty or null")
	} else {
		Log.Info().Msg("Unsubscribing from old Channels")
		for _, ct := range currentCT.Channels {
			if newCT != nil && newCT.Channels != nil {
				for _, nt := range newCT.Channels {
					if ct.LinkID == nt.LinkID {
						break // already unsubscribed
					}
				}
			}
			currentCT.topicMan.UnsubscribeTopic(ct.constructTopic()) // unsubscribe
		}
	}
	return true
}

var i int = 0

func initChannelTopics(wg *sync.WaitGroup) (*ChannelTopics, error) {
	d, _ := time.ParseDuration(viper.GetString("CONFIG_REFRESH_INTERVAL"))
	ticker := time.NewTicker(d)
	channelTopics := ChannelTopics{done: make(chan bool)}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-channelTopics.done:

				Log.Info().Msg("received done")
				return
			case t := <-ticker.C:
				Log.Debug().Msgf("ChannelTopics Refresh Tick at %v", t)
				fetchChannelTopics()

				//testing
				top := "in.id.clyszkfc70002zpa9ooq25gq1-5lef-ldkv-sP2mEg.m.batch.t.events"
				if i != 0 {
					top = fmt.Sprintf("%v%d", top, i)
				}
				i = i + 1
				channelTopics.topicMan.SubscribeTopic(top)
			}
		}
	}()
	return &channelTopics, nil
}

func (ct *ChannelTopics) AttachTopicMan(tm *TopicMan) {
	ct.topicMan = tm
}

func (ct *ChannelTopics) Close() {
	ct.done <- true
}
