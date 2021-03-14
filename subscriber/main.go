package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	windowSumTopic    = "window_sum"
	windowMedianTopic = "window_median"
	doneMsg           = "DONE" // signal from publisher to finish the task
)

var outFolder string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	outFolder = filepath.Join(cwd, "output")
	log.Printf("Output folder: %v\n", outFolder)
}

func receive(topic *redis.PubSub, subscriber Subscriber, wg *sync.WaitGroup) {
	count := 0
	for msg := range topic.Channel() {
		content := msg.Payload
		if content == doneMsg {
			break
		}
		num, err := strconv.Atoi(content)
		if err != nil {
			log.Fatalln("Unexpected content received, expect an integer.")
		} else {
			subscriber.Receive(num)
		}
		count++
	}
	wg.Done()
	log.Printf("Topic %s received %d messages in total.\n", topic.String(), count)
}

func main() {
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

	// subscribe to two topics
	ctx := context.Background()
	sumTopic := client.Subscribe(ctx, windowSumTopic)
	medianTopic := client.Subscribe(ctx, windowMedianTopic)
	defer func() {
		sumTopic.Close()
		medianTopic.Close()
	}()

	// create output files
	sumOutFile, err := os.Create(filepath.Join(outFolder, windowSumTopic+"_result.txt"))
	if err != nil {
		panic(err)
	}
	medianOutFile, err := os.Create(filepath.Join(outFolder, windowMedianTopic+"_result.txt"))
	if err != nil {
		panic(err)
	}
	defer func() {
		sumOutFile.Close()
		medianOutFile.Close()
	}()

	// add wait groups for tracking subscribers
	var wg sync.WaitGroup
	wg.Add(2)

	// do receive
	sumSub := NewSumSubscriber(sumOutFile)
	medianSub := NewMedianSubscriber(medianOutFile)
	go receive(sumTopic, sumSub, &wg)
	go receive(medianTopic, medianSub, &wg)

	wg.Wait()
}
