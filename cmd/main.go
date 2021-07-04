package main

import (
	"context"
	"log"
	"time"

	pb "github.com/myriadeinc/zircon-proto"
	"google.golang.org/grpc"
	prototext "google.golang.org/protobuf/encoding/prototext"
)

const (
	address = "localhost:8088"
)

func main() {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewZirconClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pbb := pb.Block{
		HexResult: "0c456307f6681606439dcfebfa890473d10dc85fdadf3da4909b9e1838380000",
		// HexNonce:   "0c456307f6681606439dcfebfa890473d10dc85fdadf3da4909b9e1838380000",
		GlobalDiff: "12345670",
		LocalDiff:  "10000000",
	}
	r, err := c.ProcessBlock(ctx, &pbb)
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}
	if r.GetBlockStatus() == pb.PatriciaBlockResponse_VALID {
		log.Println("valid one!")
	}
	// log.Printf("We got: %s", r.GetTest())

	log.Printf(prototext.Format(r))

}
