package tools

import (
	"github.com/spf13/viper"
	"net/http"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"time"
	"context"
	"math/rand"
	"bytes"
	"encoding/json"
)

//read the config file, helper function
func ReadConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath("./configs")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}

//general GET request
func GetRequest(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Panicf("%s", err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Panicf("%s", err)
		}
		return contents
	}
	return nil
}

//better version //TODO: remove the above version after testing
func GetRequest2(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicf("%v", err)
		return ""
	}

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Millisecond)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Panicf("%v", err)
		return ""
	}

	if b, err := ioutil.ReadAll(res.Body); err == nil {
		return string(b)
	}
	return ""
}

//general POST request
func PostJsonRequest(url string, jsonMap map[string]interface{}) ([]byte, error) {

	jsonValue, err := json.Marshal(jsonMap)
	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	log.Printf("%s we got ",string(b))

	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}

	return b, nil
}

// Generate a Random String of length n
func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// quick function to check for an error and, optionally terminate the program
func ErrorCheck(err error, where string, kill bool) {
	if err != nil {
		if kill {
			log.WithError(err).Fatalln("Script Terminated")
		} else {
			log.WithError(err).Warnf("@ %s\n", where)
		}
	}
}
