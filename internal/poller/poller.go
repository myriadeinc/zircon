package poller

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/myriadeinc/zircon/internal/cache"
	"github.com/myriadeinc/zircon/internal/nodeapi"
	"github.com/myriadeinc/zircon/internal/server"
)

type Poller struct {
	client      nodeapi.NodeApi
	cache       cache.CacheService
	blockHeight uint64
}

func NewPoller() Poller {
	nodes := []string{"https://node.monerod.org/json_rpc"}

	client := nodeapi.NewNodeClient(nodes)
	cache := cache.NewClient()

	return Poller{
		client:      client,
		cache:       cache,
		blockHeight: 0,
	}
}

func (p *Poller) PollForever() {
	log.Info().Uint64("blockheight", p.blockHeight).Msg("Starting poller")
	for {
		template := p.client.GetRawBlockTemplate()
		// Case of bad startup
		if len(template) == 0 {
			log.Error().Msg("Could not get block template from node")
			time.Sleep(15 * time.Second)
			continue
		}

		height := uint64(template["height"].(float64))

		err := p.cache.SaveNewTemplate(template)
		if err != nil {
			log.Error().Err(err).Msg("could not save newtemplate")
			time.Sleep(15 * time.Second)
			continue
		}
		log.Trace().Msg("Saved new blocktemplate")
		if height > p.blockHeight {
			log.Info().Uint64("height", height).Msg("new blockheight detected")
			p.blockHeight = height
			// Fire and forget
			sessions := server.GetSessionHandler()
			go sessions.BroadcastNewJobs()

		}

		time.Sleep(15 * time.Second)

	}

}
