/*
 * Copyright (c) 2024 valmi.io <https://github.com/valmi-io>
 *
 * Created Date: Wednesday, July 17th 2024, 6:11:58 pm
 * Author: Rajashekar Varkala @ valmi.io
 */

package util

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/spf13/viper"
	. "github.com/valmi-io/cx-pipeline/internal/log"
)
 


func SetConfigAuth(req *http.Request){
	req.SetBasicAuth(viper.GetString("APP_BACKEND_USER"), viper.GetString("APP_BACKEND_PASSWORD"))
}

func GetUrl(url string, setAuth func(*http.Request)) (string, int, error) {
	client := http.Client{Timeout: 5 * time.Second}

    req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
    if err != nil {
        Log.Error().Msg(err.Error())
		return "", -1, err

    }
    if setAuth!= nil {
        setAuth(req)
    }

	resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close() 

    Log.Info().Msgf("Response status: %v", resp.Status) 

    resBody, err := io.ReadAll(resp.Body)
    if err != nil {
        Log.Error().Msg(err.Error())
		return "", resp.StatusCode, err
    }
	return string(resBody), resp.StatusCode, nil

}

func PostUrl(url string, body []byte, setAuth func(*http.Request)) (string, int, error) {
	client := http.Client{Timeout: 5 * time.Second}

    req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
    if err != nil {
        Log.Error().Msg(err.Error())
		return "", -1, err

    }
    if setAuth!= nil {
        setAuth(req)
    }

	resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close() 

    Log.Info().Msgf("Response status: %v", resp.Status) 

    resBody, err := io.ReadAll(resp.Body)
    if err != nil {
        Log.Error().Msg(err.Error())
		return "", resp.StatusCode, err
    }
	return string(resBody), resp.StatusCode, nil

}