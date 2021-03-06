package tools

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"bufio"
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
func GETRequest(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warnf("%v", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(req.Context(), 100*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Warnf("%v", err)
		return nil
	}

	if contents, err := ioutil.ReadAll(res.Body); err == nil {
		return contents
	}
	return nil
}

//general POST request
func POSTRequest(url string, payload []byte) ([]byte, error) {

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}

	if contents, err := ioutil.ReadAll(resp.Body); err == nil {
		log.Info("POST Request Returned >>> " + string(contents))
		return contents, nil
	}
	return nil, err
}

// general PUT request
func PUTRequest(url string, payload []byte) ([]byte, error) {

	body := bytes.NewReader(payload)

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if contents, err := ioutil.ReadAll(resp.Body); err == nil {
		log.Info("PUT Request Returned >>> " + string(contents))
		return contents, nil
	}
	return nil, err
}

// general DELETE request
func DELETERequest(url string) ([]byte, error) {

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if contents, err := ioutil.ReadAll(resp.Body); err == nil {
		log.Info("DELETE Request Returned >>> " + string(contents))
		return contents, nil
	}
	return nil, err

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

//convert hex to int
func HexToUInt(hexStr string) uint64 {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint64(result)
}

//generate sha1 hash from interface{}
func GetSha1Hash(payload interface{}) string {

	out, err := json.Marshal(payload)
	if err != nil {
		log.Error(err)
		return ""
	}

	algorithm := sha1.New()
	algorithm.Write(out)
	return fmt.Sprintf("%x", algorithm.Sum(nil))
}

// google-chrome-stable needs to be installed
func GeneratePdf(fromFile string, toFile string) error {

	//test if google-chrome-stable is installed

	cmd := exec.Command("/usr/bin/google-chrome-stable", "-version")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Printf("everything looks good -> %s", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	// ======= CORE =========
	cmd = exec.Command("/usr/bin/google-chrome-stable", "--headless", "--disable-gpu", "--virtual-time-budget=1000",
		"--print-to-pdf=/home/ubuntu/go/src/github.com/motionwerkGmbH/cpo-backend-api/"+toFile, "/home/ubuntu/go/src/github.com/motionwerkGmbH/cpo-backend-api/"+fromFile)
	cmdReader, err = cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner = bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Printf("%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

// return the external ip
func GetExternalIp() []byte {
	req, err := http.NewRequest("GET", "http://ipecho.net/plain", nil)
	if err != nil {
		log.Errorf("%v", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(req.Context(), 20*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Errorf("%v", err)
		return nil
	}

	if contents, err := ioutil.ReadAll(res.Body); err == nil {
		return contents
	}
	return nil
}
