package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

// upper limit of the random integer sent to subscriber
// To avoid integer overflow, make sure 20RPS * 5sec * maxN < MaxInt64
// To make numbers diverse enough (the size of queue will be up to 100), set maxN = 10000
const maxN = 10000
const doneMsg = "DONE"

var sendDuration time.Duration

func init() {
	// TODO: read from config
	sendDuration = 60 * time.Second
}

func main() {
	// rand.Seed(0)

	// connect to redis
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()

	time.Sleep(1 * time.Second)
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Println("Cannot connect to Redis, will retry once...")
		time.Sleep(2 * time.Second)

		if _, err = client.Ping(context.Background()).Result(); err != nil {
			log.Println("Failed to connect to Redis")
			panic(err)
		}
	}

	// started to send for ${sendDuration} seconds
	log.Printf("Started to send random integers to subscriber for %v\n", sendDuration)

	ctx := context.Background()
	t := time.Now()
	for i := 1; ; i++ {
		now := time.Now()
		elasped := now.Sub(t)

		if elasped > sendDuration {
			// stop sending to subscribers.
			// a workaround to send a singal to all subscribers for them to unsubscribe
			go func() {
				if err := client.Publish(ctx, "window_sum", doneMsg).Err(); err != nil {
					panic(err)
				}
				if err := client.Publish(ctx, "window_median", doneMsg).Err(); err != nil {
					panic(err)
				}
			}()

			log.Printf("Sent %d messages in total, elasped time %v.\n", i-1, elasped)
			break
		}

		// choose a random number to send
		n := rand.Intn(maxN)
		go func() {
			if err := client.Publish(ctx, "window_sum", n).Err(); err != nil {
				panic(err)
			}
			if err := client.Publish(ctx, "window_median", n).Err(); err != nil {
				panic(err)
			}
		}()

		if i%100 == 0 {
			log.Printf("Sent %d messages, elasped time %v.\n", i, elasped)
		}

		time.Sleep((50 * time.Millisecond) - time.Now().Sub(now)) // 20 RPS
	}

	time.Sleep(3 * time.Second)
}
