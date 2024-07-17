/*
 * Copyright (c) 2024 valmi.io <https://github.com/valmi-io>
 *
 * Created Date: Wednesday, July 17th 2024, 6:11:58 pm
 * Author: Rajashekar Varkala @ valmi.io
 */

package configstore

type ConfigStore struct {
	channelTopics *ChannelTopics
	storefrontIfttts *StorefrontIfttts
}

func Init() (*ConfigStore, error) {
	channelTopics, error := initChannelTopics()
	if error!= nil {
        return nil, error
    }

	storefrontIfttts, error := initStoreFrontIfttts()
	if error!= nil {
        return nil, error
    }

	return &ConfigStore{
        channelTopics: channelTopics,
		storefrontIfttts: storefrontIfttts,
    }, nil
}

func (cs *ConfigStore) Close(){
	cs.channelTopics.Close()
    cs.storefrontIfttts.Close()
}