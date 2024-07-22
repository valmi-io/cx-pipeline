/*
 * Copyright (c) 2024 valmi.io <https://github.com/valmi-io>
 *
 * Created Date: Wednesday, July 17th 2024, 6:11:58 pm
 * Author: Rajashekar Varkala @ valmi.io
 */

package configstore

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/spf13/viper"
	. "github.com/valmi-io/cx-pipeline/internal/log"
	util "github.com/valmi-io/cx-pipeline/internal/util"
)

func fetchIfttts() (StorefrontIfttts, error) {
	data, respCode, err := util.GetUrl(viper.GetString("APP_BACKEND_URL")+"/api/v1/superuser/ifttts", util.SetConfigAuth)
	Log.Info().Msg(data)
	Log.Info().Msgf("%v", respCode)
	if err != nil {
		Log.Info().Msg(err.Error())
	}

	var storeIfttts []StoreIfttt
	unmarshallErr := json.Unmarshal([]byte(data), &storeIfttts)
	if unmarshallErr != nil {
		Log.Error().Msgf("Error Unmarshalling JSON: %v", unmarshallErr)
		return StorefrontIfttts{}, unmarshallErr
	}

	return StorefrontIfttts{done: make(chan bool), StoreIfttts: storeIfttts}, nil
}

type StoreIfttt struct {
	StoreID string `json:"store_id"`
	Code    string `json:"code"`
}

type StorefrontIfttts struct {
	StoreIfttts []StoreIfttt
	done        chan bool
}

func initStoreFrontIfttts(wg *sync.WaitGroup) (*StorefrontIfttts, error) {
	d, _ := time.ParseDuration(viper.GetString("CONFIG_REFRESH_INTERVAL"))
	ticker := time.NewTicker(d)
	storefrontIfttts := StorefrontIfttts{done: make(chan bool)}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-storefrontIfttts.done:
				Log.Info().Msg("received done")
				ticker.Stop()
				return
			case t := <-ticker.C:
				Log.Debug().Msgf("StorefrontIfttts Refresh Tick at %v", t)
				newStorefrontIfttts, err := fetchIfttts()
				if err != nil {
					Log.Error().Msgf("Error fetching ifttts: %v", err)
					continue
				}
				storefrontIfttts = newStorefrontIfttts
			}
		}
	}()
	return &storefrontIfttts, nil
}

func (si *StorefrontIfttts) Close() {
	si.done <- true
}
