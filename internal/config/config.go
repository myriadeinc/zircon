package config

import (
	"github.com/spf13/viper"
)

type ConfigKeys struct {
}

func DefaultConfigs() {
	viper.Set("WS_RPC_URL", "ws://0.0.0.0:9877")

}