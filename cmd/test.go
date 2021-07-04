package main

import (
	"log"

	pb "github.com/myriadeinc/zircon_proto"
)

func main() {

	var s = pb.TestLib()

	log.Printf(s)
}
