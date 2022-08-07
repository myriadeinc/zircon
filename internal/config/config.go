package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type ConfigKeys struct {
}

func DefaultConfigs() {
	viper.SetDefault("WS_RPC_URL", "ws://0.0.0.0:9877")
	viper.SetDefault("TRACE_LOGS", false)
	viper.SetDefault("WALLET_ADDRESS", "47PAULmUFo3DPHKehGPuxXbEAB4JkRYJ49DEFs4EqaT7M2TRqqWWHAeJyEHWg8eRoWNwMAHh7bx6Eh5SR2fpdnj71fhxugC")
	viper.SetDefault("RESERVE_OFFSET", 8)
	viper.SetDefault("REDIS_URL", "redis://cache")
	viper.SetDefault("FORCE_JOB_PUSH", false)
	viper.SetDefault("MONERO_NODES", "https://node.monerod.org/json_rpc,https://monero.2255.me,https://chad.fiatfaucet.com/json_rpc,http://82.65.156.176:18081/json_rpc")
}

func DumpConfigs() {
	log.Debug().Str("WS_RPC_URL", viper.GetString("WS_RPC_URL")).Msg("")
	log.Debug().Str("REDIS_URL", viper.GetString("REDIS_URL")).Msg("")
	log.Debug().Str("WALLET_ADDRESS", viper.GetString("WALLET_ADDRESS")).Msg("")
	log.Debug().Str("MONERO_NODES", viper.GetString("MONERO_NODES")).Msg("")
	log.Debug().Int("RESERVE_OFFSET", viper.GetInt("RESERVE_OFFSET")).Msg("")
	log.Debug().Bool("TRACE_LOGS", viper.GetBool("TRACE_LOGS")).Msg("")
	log.Debug().Bool("FORCE_JOB_PUSH", viper.GetBool("FORCE_JOB_PUSH")).Msg("")

}
