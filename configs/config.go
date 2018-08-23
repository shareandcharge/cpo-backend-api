package configs

import (
	log "github.com/Sirupsen/logrus"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"github.com/spf13/viper"
)

func Load() *viper.Viper {
	// Configs
	Config, err := tools.ReadConfig("api_config", map[string]interface{}{
		"port":        9090,
		"hostname":    "localhost",
		"environment": "debug",
		"cpo": map[string]string{
			"wallet_address": "0x5c9b043d100a8947e614bbfdd8c6077bc7c456d0",
			"wallet_seed":    "distance hover flock tomorrow fault rain decline magic teach impact cart drum",
		},
	})
	if err != nil {
		log.Error("Error when reading config: %v\n", err)
	}
	return Config
}
