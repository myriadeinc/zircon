package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/myriadeinc/zircon/internal/models"
	"github.com/spf13/viper"

	"github.com/go-redis/redis/v8"
)

type CacheService interface {
	SaveNewTemplate(models.StrictTemplate) error
	FetchTemplate() (*models.StrictTemplate, error)
}

type RedisService struct {
	client       *redis.Client
	minerRecords map[string]struct{}
}

func NewClient() CacheService {

	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisService{
		client:       rdb,
		minerRecords: map[string]struct{}{},
	}

}

const key = "blocktemplate"

func (r *RedisService) SaveNewTemplate(template models.StrictTemplate) error {
	return r.client.Set(context.Background(), key, template, 0).Err()
}

func (r *RedisService) FetchTemplate() (*models.StrictTemplate, error) {
	jsonString, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	var template models.StrictTemplate
	err = json.Unmarshal([]byte(jsonString), &template)
	if err != nil {
		return nil, err
	}
	if !template.IsValid() {
		return nil, fmt.Errorf("expected full template, received %v", template)
	}
	return &template, nil
}

func (r *RedisService) UpsertMinerDifficulty(minerId string, blockHeight string, difficulty uint64) error {
	r.minerRecords[minerId] = struct{}{}
	key := fmt.Sprintf("%s_%s", blockHeight, minerId)
	return r.client.Set(context.Background(), key, difficulty, 0).Err()
}

func (r *RedisService) clearMiners() {
	r.minerRecords = map[string]struct{}{}
}

func (r *RedisService) FlushMiners(blockHeight string) error {
	defer r.clearMiners()
	for minerId := range r.minerRecords {
		key := fmt.Sprintf("%s_%s", blockHeight, minerId)
		_, err := r.client.Get(context.Background(), key).Result()
		return err
	}

	return nil

}
