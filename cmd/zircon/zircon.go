package main

import (
	// "context"
	// "log"
	// "time"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	server "github.com/myriadeinc/zircon/internal/server"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	server := server.New()
	server.Listen("0.0.0.0:12345")

}
