package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"github.com/go-redis/redis/v8"
)

type CacheService interface {
	SaveNewTemplate(map[string]interface{}) error
	FetchTemplate() (*StrictTemplate, error)
}

// Instead of map[string]string because the redis library only has support to serialize a struct
type StrictTemplate struct {
	BlockTemplateBlob string `json:"blocktemplate_blob"`
	Difficulty        string `json:"difficulty"`
	SeedHash          string `json:"seed_hash"`
	Height            string `json:"height"`
}

type RedisService struct {
	client *redis.Client
}

func NewClient() CacheService {

	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisService{
		client: rdb,
	}

}

const key = "blocktemplate"

func (r *RedisService) SaveNewTemplate(template map[string]interface{}) error {

	redisTemplate := StrictTemplate{
		BlockTemplateBlob: fmt.Sprintf("%v", template["blocktemplate_blob"]),
		SeedHash:          fmt.Sprintf("%v", template["seed_hash"]),
		// Horrible yes, but desired as they are parsed as floats but should be uint
		Difficulty: fmt.Sprintf("%d", uint64(template["difficulty"].(float64))),
		Height:     fmt.Sprintf("%d", uint64(template["height"].(float64))),
	}

	return r.client.Set(context.Background(), key, redisTemplate, 0).Err()
}

func (r *RedisService) FetchTemplate() (*StrictTemplate, error) {
	jsonString, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	var template StrictTemplate
	err = json.Unmarshal([]byte(jsonString), &template)
	if err != nil {
		return nil, err
	}
	if len(template.BlockTemplateBlob) == 0 {
		return nil, errors.New("template has no blob")
	}
	return &template, nil
}

func (s StrictTemplate) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
