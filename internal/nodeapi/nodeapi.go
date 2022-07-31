package nodeapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type NodeApi interface {
	GetRawBlockTemplate() map[string]interface{}
}

type NodeClient struct {
	nodes []string
}

func NewNodeClient(nodes []string) NodeApi {
	client := NodeClient{
		nodes: nodes,
	}
	return &client
}

// Block template has these fields
// {

// 	"blockhashing_blob": "0e0e8c89eb9606b44e67bdbc91e5f95c5ea49ba410b1c223ab897417b19a0f631b80cf5be43a4e00000000f608c45071e983a7fbcd9d1474bfa1505d4d808b2bfcd513f214a3820c41942d17",
// 	"blocktemplate_blob": "0e0e8c89eb9606b44e67bdbc91e5f95c5ea49ba410b1c223ab897417b19a0f631b80cf5be43a4e0000000002fc92a30101ffc092a30101f0e49fedbc11027d556d1fc303296a55dec60ca1394c7f49ecb5b06a9c26d9dfc304d91f4cf9752b014283b855539052e430134848504dfb8632191849f72bed650c4abf7adb106bbd020800000000000000000016a2ef6ca421eba119678fbbed31ad7a8af2952fe783d98f1f1f4b3f7ec9969bbf56d627b2d0c178273bc68fa27c30663f67036b13afe46a2709317a209b3f1f17d324fb6cb441322a51dece5981b3b160784cfa1e38f366ffef06aed8a9605a06fc0a5431b8021d18bd74f847237f7fe675df131eaf94cb17393ba3e7405d2881ed7d1ef26878e861848558d9da280d0cc5014ae7e154ae12c2d542f65b5d270bd2b16841221c6502e514d61ca55a178e170f618c63124ac804d6f10a6cd8baa075b07558c83eba5e986e4d6f8a95386c52e4d4777a650b0d43052fe5859b4528cd19acc88220fbfdfb38434cf06b21b4969ddf4554de500964720eae235deedc575f9804e293faf1b1484890a807cb950df2d40e475e170cf488719bb45416475d239e64b22502a045849cc3c37859405a0630590b35e63e4495b24a2b950af1edd685ba86e04184b05555b3dc2c57f0c566016110ad551b9b6a5a62c08071288b1eaf4edb80cbe474ec13eb3836f92a047f42e2098a639c1da6fbf5ba1c4fba12992e6676effc8178a71c833022362aee161f334530f917be43041b92b26bdc12356b83c44059413425e34d6a8f09ac1ea34fc6ce5b88df15e368e3cf3da93396e4eded2501c8aa1957e991bd6de0e13309b1521364c730007fba813b41c46c5c7f3cb2088ea6ea799f508f887f603b67d87d4afc36b7cce441da21ebeb3b22ea07ca4e97284ba0548af9a02ea2dce5b709a3a2c5b90435674845a9aad23928829e9e20eaaf42ce612e9b5f0fffb79670ecee7cfcc29bc9d3fe7f47840849b482f46c0e6497e9f086574ccf84d609bda28ce0ee8928ef00eef0cbc07408370dd4f6418c17e9badd0230d384aea70f92102b7ebf4f3bac414130f5e853753df403f80c96b962a6c54f8254069c8cfdcc8d6d4d66023cf2102beeea07ba7a74d0c01ce4658f33059c2c2989e7ed8fa8136160d8d1b163ae22692e95c26a988be8",
// 	"difficulty": 307042777507,
// 	// "difficulty": 77507,
// 	"difficulty_top64": 0,
// 	"expected_reward": 600450790000,
// 	"height": 2672960,
// 	"next_seed_hash": "",
// 	"prev_hash": "b44e67bdbc91e5f95c5ea49ba410b1c223ab897417b19a0f631b80cf5be43a4e",
// 	"reserved_offset": 130,
// 	"seed_hash": "348b6ead39ba7ebb7ab45bbada3e8e7c3c171e257671e5ac570867c953d5bcf9",
// 	"seed_height": 2672640,
// 	"status": "OK",
// 	"untrusted": false,
// 	"wide_difficulty": "0x477d2cf9a3"

// }
func (n *NodeClient) GetRawBlockTemplate() map[string]interface{} {
	walletAddress := viper.GetString("WALLET_ADDRESS")
	reserveOffset := viper.GetInt("RESERVE_OFFSET")

	request, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "get_block_template",
		"params": map[string]interface{}{
			"wallet_address": walletAddress,
			"reserve_size":   reserveOffset,
		},
	})

	requestBody := bytes.NewBuffer(request)

	for _, node := range n.nodes {
		resp, err := http.Post(node, "application/json", requestBody)
		if err != nil {
			log.Error().Err(err).Str("url", node).Msg("could not use node")
			resp.Body.Close()
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Str("url", node).Msg("receive bad response")
			resp.Body.Close()
			continue
		}
		return convertResponseToTemplate(body)
	}
	return nil
}

func convertResponseToTemplate(response []byte) map[string]interface{} {
	holder := struct {
		Result map[string]interface{} `json:"result,omitempty"`
	}{}

	err := json.Unmarshal(response, &holder)
	if err != nil {
		log.Error().Err(err).Msg("unmarshal failed for rpc response")
		return map[string]interface{}{}
	}
	return holder.Result

}
