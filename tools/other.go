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

// a general get request with 100 seconds timeout
func GetRequest(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicf("%v", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(req.Context(), 100*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Panicf("%v", err)
		return nil
	}

	if contents, err := ioutil.ReadAll(res.Body); err == nil {
		return contents
	}
	return nil
}

//general POST request
func PostRequest(url string, payload []byte) ([]byte, error) {


	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
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

// general PUT request
func PUTRequest(url string, payload []byte) ([]byte, error) {

	req, err := http.NewRequest(http.MethodPut, url,  bytes.NewBuffer(payload))
	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}

	b, err := ioutil.ReadAll(req.Body)

	log.Printf("%s we got ",string(b))

	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}

	return b, nil
}

// general DELETE request
func DELETERequest(url string) ([]byte, error) {

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()


	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("%s we got ",string(b))
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
