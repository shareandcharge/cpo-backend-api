package configs

import (
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"fmt"
	"github.com/spf13/viper"
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
