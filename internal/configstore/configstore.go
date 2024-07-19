/*
 * Copyright (c) 2024 valmi.io <https://github.com/valmi-io>
 *
 * Created Date: Wednesday, July 17th 2024, 6:11:58 pm
 * Author: Rajashekar Varkala @ valmi.io
 */

package configstore

import (
	"sync"

	. "github.com/valmi-io/cx-pipeline/internal/msgbroker"
)

type ConfigStore struct {
	channelTopics    *ChannelTopics
	storefrontIfttts *StorefrontIfttts
	wg               *sync.WaitGroup
}

func Init() (*ConfigStore, error) {
	var wg sync.WaitGroup
	channelTopics, error := initChannelTopics(&wg)
	if error != nil {
		return nil, error
	}

	storefrontIfttts, error := initStoreFrontIfttts(&wg)
	if error != nil {
		return nil, error
	}

	return &ConfigStore{
		channelTopics:    channelTopics,
		storefrontIfttts: storefrontIfttts,
		wg:               &wg,
	}, nil
}

func (cs *ConfigStore) AttachTopicMan(tm *TopicMan) {
	cs.channelTopics.topicMan = tm
}

func (cs *ConfigStore) Close() {
	cs.channelTopics.Close()
	cs.storefrontIfttts.Close()
	cs.wg.Wait()
}
