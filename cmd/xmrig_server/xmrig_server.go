package main

import (
	// "context"
	// "log"
	// "time"

	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	config "github.com/myriadeinc/zircon/internal/config"
	server "github.com/myriadeinc/zircon/internal/server"
	"github.com/spf13/viper"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	config.DefaultConfigs()
	viper.AutomaticEnv()

	useTraceLevel := viper.GetBool("TRACE_LOGS")
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if useTraceLevel {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	config.DumpConfigs()

	pool := server.New()

	// We don't really need to do this because we do reconnects, but still helpful
	time.Sleep(10 * time.Second)

	pool.Start("0.0.0.0:8222")

}
