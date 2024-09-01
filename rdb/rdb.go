package rdb

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-mailer"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx context.Context

func Init() {
	ctx = context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Instance() *redis.Client {
	return rdb
}

func PublishEmail(channel string, email mailer.RedisQueueEmail) error {
	payload, err := json.Marshal(email)

	if err != nil {
		return err
	}

	return Publish("email", payload)
}

func Publish(channel string, data []byte) error {
	return rdb.Publish(ctx, channel, data).Err()
}
