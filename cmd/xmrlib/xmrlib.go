package main

import (
	// "context"
	// "log"
	// "time"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	xmrlib "github.com/myriadeinc/zircon/xmrlib"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msgf("%v", xmrlib.Hello())

}
