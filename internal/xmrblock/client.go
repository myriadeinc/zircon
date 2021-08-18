package xmrblock

import (
	"encoding/json"
	"fmt"

	"github.com/ybbus/jsonrpc"
)

var client jsonrpc.RPCClient

type BlockTemplateRequest struct {
	Wallet string `json:"wallet_address"`
	Size   uint   `json:"reserve_size"`
}

func InitClient() {
	// client = jsonrpc.NewClient("http://node.melo.tools:18081")
	// client = jsonrpc.NewClient("http://daemon.myriade.io")
}
func GetDaemonBlockTemplate() *json.RawMessage {
	wallet := "42PooLTYHzzZPY15hZt5SJVRCLYZCLYPdMhRjXaUY3SwERq9yxCnKg1EnGd1YqTR9GfH1EK4LfquBbRPuWKtPku6M6qhyZg"
	client = jsonrpc.NewClient("http://node.melo.tools:18081/json_rpc")

	var response *json.RawMessage
	err := client.CallFor(&response, "get_block_template", &BlockTemplateRequest{
		Wallet: wallet,
		Size:   8,
	})
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
	}
	return response
}
