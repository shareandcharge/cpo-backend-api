package configs

import (
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"encoding/json"
	"io/ioutil"
	"log"
)

func Load() (*viper.Viper) {
	// Configs
	Config, err := tools.ReadConfig("api_config", map[string]interface{}{
		"port":     9090,
		"hostname": "localhost",
		"auth": map[string]string{
			"username": "user",
			"password": "pass",
		},
	})
	if err != nil {
		panic(fmt.Errorf("Error when reading config: %v\n", err))
	}
	return Config
}

//updats the seed in ~/.sharecharge/config.json

func UpdateBaseAccountSeedInSCConfig(seed string){

	type ConfigStruct struct {
		TokenAddress  string `json:"tokenAddress"`
		LocationsPath string `json:"locationsPath"`
		TariffsPath   string `json:"tariffsPath"`
		BridgePath    string `json:"bridgePath"`
		Seed          string `json:"seed"`
		Stage         string `json:"stage"`
		GasPrice      int    `json:"gasPrice"`
		EthProvider   string `json:"ethProvider"`
		IpfsProvider  struct {
			Host     string `json:"host"`
			Port     string `json:"port"`
			Protocol string `json:"protocol"`
		} `json:"ipfsProvider"`
	}

	jsonFile, err := os.Open("~/.sharecharge/config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	config := ConfigStruct{}
	err = json.Unmarshal(byteValue, &config)
	tools.ErrorCheck(err, "config.go", false)


	fmt.Println("Successfully Opened config.json")
	log.Printf("%s", config)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()


}