package main

import (
	// "context"
	// "log"
	// "time"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	server "github.com/myriadeinc/zircon/internal/server"
	// zirconbuf "github.com/myriadeinc/zircon_proto"
	// "google.golang.org/grpc"
	// prototext "google.golang.org/protobuf/encoding/prototext"
	xmrlib "github.com/myriadeinc/zircon/xmrlib"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msgf("%v", xmrlib.Hello())

	server := server.GetServerInstance()
	server.Listen("0.0.0.0:5656")

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	// zirconbufb := zirconbuf.Block{
	// 	HexResult: "0c456307f6681606439dcfebfa890473d10dc85fdadf3da4909b9e1838380000",
	// 	// HexNonce:   "0c456307f6681606439dcfebfa890473d10dc85fdadf3da4909b9e1838380000",
	// 	GlobalDiff: "12345670",
	// 	LocalDiff:  "10000000",
	// }

	// r, err := c.ValidateBlock(ctx, &zirconbufb)
	// if err != nil {
	// 	log.Fatalf("could not call: %v", err)
	// }
	// if r.GetBlockStatus() == zirconbuf.BlockResponse_VALID {
	// 	log.Println("valid one!")
	// }

	// log.Printf(prototext.Format(r))

}
