package configstore

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	. "github.com/valmi-io/cx-pipeline/internal/log"
	. "github.com/valmi-io/cx-pipeline/internal/msgbroker"
	util "github.com/valmi-io/cx-pipeline/internal/util"
)

type ChannelTopic struct {
	LinkID     string `json:"link_id"`
	WriteKey   string `json:"write_key"`
	storefront string
	channel    string
}

type ChannelTopics struct {
	mu       sync.RWMutex
	Channels []ChannelTopic
	done     chan bool
	topicMan *TopicMan
}

func fetchChannelTopics(currentCT *ChannelTopics) []ChannelTopic {
	jsonPayload := `{"channel_in": ["chatbox"], "channel_not_in": ["x", "y"]}`
	data, respCode, err := util.PostUrl(
		viper.GetString("APP_BACKEND_URL")+"/api/v1/superuser/channeltopics",
		[]byte(jsonPayload),
		util.SetConfigAuth,
		nil)

	Log.Info().Msg(data)
	Log.Info().Msgf("%v", respCode)
	if err != nil {
		Log.Info().Msg(err.Error())
	}
	// PARSE the response and store it in the ChannelTopics struct
	var newChannelTopics []ChannelTopic
	if err := json.Unmarshal([]byte(data), &newChannelTopics); err != nil {
		Log.Error().Msgf("Error Unmarshalling JSON: %v", err)
	}

	// Pass to matchChannelState
	if status := matchChannelState(newChannelTopics, currentCT.Channels, currentCT.topicMan); status {
		// if success, switch to NewChannelTopics
		return newChannelTopics
	}
	return currentCT.Channels
}

func (t *ChannelTopic) constructTopic() string {
	return fmt.Sprintf("in.id.%v.m.batch.t.events", t.LinkID)
}

func matchChannelState(newCT []ChannelTopic, currentCT []ChannelTopic, topicMan *TopicMan) bool {
	if topicMan == nil {
		return false
	}
	Log.Info().Msg("Subscribing to new Channels")
	newCtMap := make(map[string]ChannelTopic)
	for _, channelTopic := range newCT {
		newCtMap[channelTopic.LinkID] = channelTopic
	}
	currentCTmap := make(map[string]ChannelTopic)
	for _, channelTopic := range currentCT {
		currentCTmap[channelTopic.LinkID] = channelTopic
	}
	for k, v := range newCtMap {
		if _, found := currentCTmap[k]; !found {
			Log.Info().Msgf("supradeep: %v", v.constructTopic())
			topicMan.SubscribeTopic(v.constructTopic())
		}
	}
	for k, v := range currentCTmap {
		if _, found := newCtMap[k]; !found {
			Log.Info().Msgf("UnSubscribing to old Channel: %v", v.constructTopic())
			topicMan.UnsubscribeTopic(v.constructTopic()) //Unsubcrbe
		}
	}
	return true
}

// var i int = 0

func initChannelTopics(wg *sync.WaitGroup) (*ChannelTopics, error) {
	d, _ := time.ParseDuration(viper.GetString("CONFIG_REFRESH_INTERVAL"))
	ticker := time.NewTicker(d)
	channelTopics := &ChannelTopics{done: make(chan bool)}
	firstTime := true
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-channelTopics.done:

				Log.Info().Msg("received done")
				return
			case t := <-ticker.C:
				if !firstTime {
					continue
				}
				Log.Debug().Msgf("ChannelTopics Refresh Tick at %v", t)
				newChannels := fetchChannelTopics(channelTopics)
				channelTopics.mu.Lock()
				channelTopics.Channels = newChannels
				channelTopics.mu.Unlock()
				Log.Debug().Msgf("ChannelTopics Refresh Tick at %+v", channelTopics)
				firstTime = false
				//testing
				// top := "in.id.clyszkfc70002zpa9ooq25gq1-5lef-ldkv-sP2mEg.m.batch.t.events"
				// if i != 0 {
				// 	top = fmt.Sprintf("%v%d", top, i)
				// }
				// i = i + 1
				// channelTopics.topicMan.SubscribeTopic(top)
			}
		}
	}()
	return channelTopics, nil
}

func (ct *ChannelTopics) AttachTopicMan(tm *TopicMan) {
	ct.topicMan = tm
}

func (ct *ChannelTopics) Close() {
	ct.done <- true
}
