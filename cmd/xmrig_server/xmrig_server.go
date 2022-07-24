package main

import (
	// "context"
	// "log"
	// "time"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	config "github.com/myriadeinc/zircon/internal/config"
	server "github.com/myriadeinc/zircon/internal/server"
	"github.com/spf13/viper"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	config.DefaultConfigs()
	viper.AutomaticEnv()

	pool := server.New()
	go func() {
		server.ListenWebhook()
	}()

	pool.Listen("0.0.0.0:8222")

}
