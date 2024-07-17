/*
 * Copyright (c) 2024 valmi.io <https://github.com/valmi-io>
 *
 * Created Date: Wednesday, July 17th 2024, 6:11:58 pm
 * Author: Rajashekar Varkala @ valmi.io
 */

package configstore

import (
	"time"

	"github.com/spf13/viper"
	. "github.com/valmi-io/cx-pipeline/internal/log"
	util "github.com/valmi-io/cx-pipeline/internal/util"
)

func fetchIfttts() (string, error) {
	data, respCode, err := util.GetUrl(viper.GetString("APP_BACKEND_URL") + "/api/v1/superuser/ifttts", util.SetConfigAuth)
	Log.Info().Msg(data)
	Log.Info().Msgf("%v", respCode)
	if err!= nil {
		Log.Info().Msg(err.Error())
    }
	return data, err

	// PARSE the response and store it in the StorefrontIftts struct
}

type StorefrontIfttts struct {
	StoreIfttt []struct {
        StoreID int    `json:"store_id"`
        Code   string `json:"code"`
    }
	done chan bool
}

func initStoreFrontIfttts() (*StorefrontIfttts, error) {
	d,_:= time.ParseDuration(viper.GetString("CONFIG_REFRESH_INTERVAL"))
	ticker := time.NewTicker(d)
	storefrontIfttts := StorefrontIfttts{done: make(chan bool)}
	go func() {
		for {
            select {
            case <- storefrontIfttts.done:
                return
            case t := <-ticker.C:
				Log.Debug().Msgf("StorefrontIfttts Refresh Tick at %v", t)
				fetchIfttts()
            }
        }
	}()
	return &storefrontIfttts, nil
}

func (si *StorefrontIfttts) Close(){
	si.done <- true
}


