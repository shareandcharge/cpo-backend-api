package configs

import (
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"encoding/json"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
)

func Load() (*viper.Viper) {
	// Configs
	Config, err := tools.ReadConfig("api_config", map[string]interface{}{
		"port":     9090,
		"hostname": "localhost",
		"environment": "debug",
		"msp": map[string]string{
			"wallet_address": "0xf60b71a4d360a42ec9d4e7977d8d9928fd7c8365",
			"wallet_seed": "moon another kind random mask like swarm type ostrich amused rice castle",
		},
		"cpo": map[string]string{
			"wallet_address": "0x5c9b043d100a8947e614bbfdd8c6077bc7c456d0",
			"wallet_seed": "distance hover flock tomorrow fault rain decline magic teach impact cart drum",
		},
	})
	if err != nil {
		panic(fmt.Errorf("Error when reading config: %v\n", err))
	}
	return Config
}


//updats the seed in ~/.sharecharge/config.json. Attention, the username is ubuntu. This will not work locally, unless have linux & the username "Ubuntu" :)
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

	//load the config file
	jsonFile, err := os.Open("/home/ubuntu/.sharecharge/config.json")
	tools.ErrorCheck(err, "config.go", false)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	log.Printf("%s", byteValue)
	defer jsonFile.Close()


	config := ConfigStruct{}
	err = json.Unmarshal(byteValue, &config)
	tools.ErrorCheck(err, "config.go", false)

	// this is the core of the function
	config.Seed = seed

	newconfigBytes, err := json.Marshal(config)
	tools.ErrorCheck(err, "config.go", false)


	err = ioutil.WriteFile("/home/ubuntu/.sharecharge/config.json", newconfigBytes, 644)
	tools.ErrorCheck(err, "config.go", false)

	log.Println("Successfully updated the /home/ubuntu/.sharecharge/config.json")




}