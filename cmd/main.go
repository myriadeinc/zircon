package main

import (
	"context"
	"log"
	"time"

	zirconbuf "github.com/myriadeinc/zircon_proto"
	"google.golang.org/grpc"
	prototext "google.golang.org/protobuf/encoding/prototext"
)

const (
	address = "zircon_proto:8088"
)

func main() {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := zirconbuf.NewZirconClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	zirconbufb := zirconbuf.Block{
		HexResult: "0c456307f6681606439dcfebfa890473d10dc85fdadf3da4909b9e1838380000",
		// HexNonce:   "0c456307f6681606439dcfebfa890473d10dc85fdadf3da4909b9e1838380000",
		GlobalDiff: "12345670",
		LocalDiff:  "10000000",
	}

	r, err := c.ValidateBlock(ctx, &zirconbufb)
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}
	if r.GetBlockStatus() == zirconbuf.BlockResponse_VALID {
		log.Println("valid one!")
	}

	log.Printf(prototext.Format(r))

}
