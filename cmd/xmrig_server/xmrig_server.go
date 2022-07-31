package main

import (
	// "context"
	// "log"
	// "time"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	config "github.com/myriadeinc/zircon/internal/config"
	poller "github.com/myriadeinc/zircon/internal/poller"
	server "github.com/myriadeinc/zircon/internal/server"
	"github.com/spf13/viper"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	config.DefaultConfigs()
	viper.AutomaticEnv()

	useTraceLevel := viper.GetBool("TRACE_LOGS")
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if useTraceLevel {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	config.DumpConfigs()

	pool := server.New()
	p := poller.NewPoller()
	go func() {
		p.PollForever()

	}()

	pool.Listen("0.0.0.0:8222")

}
