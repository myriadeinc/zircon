package client

// import (
// 	"fmt"
// 	"reflect"
// 	"google.golang.org/grpc"
// 	zirconbuf "github.com/myriadeinc/zircon_proto"
// )
// const (
// 	address = "zircon_proto:8088"
// )

// func InitZirconClient() error {
// 	fmt.Println("starting client")
// 	// Set up a connection to the server.
// 	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
// 	if err != nil {
// 		fmt.Println("error detected")
// 		return err
// 	}
// 	fmt.Println("connection established")
// 	defer conn.Close()
// 	c := zirconbuf.NewZirconClient(conn)
// 	fmt.Println("created new client, what type is it?")
// 	fmt.Println(reflect.TypeOf(c))
// 	return nil
// }
