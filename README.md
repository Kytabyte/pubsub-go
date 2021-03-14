# pubsub-go
A publish-subscribe system using Go-Redis

## Task

Write a publisher to send random integer with 20 RPS to two subscribers  
Subscriber 1 receives the integer and print out a sliding window sum during the past 5 seconds.  
Subscriber 2 receives the integer and print out a sliding window median during the past 5 seconds.  

## Build and Run

```bash
docker-compose up --build
```
Note: add `sudo` if needed as per your docker config

Publisher will automatically run 60 seconds, and it will stop and tell subscribers to stop.  
Output will be in `output/window_[sum|median]_result.txt`

Once the following log is presented, we can stop running the container by
pressing Ctrl+C at any time.

```
pubsub-go_subscriber_1 exited with code 0
pubsub-go_publisher_1 exited with code 0
```

## Validation

Run the following commands to run validation script

```bash
cd output
go build -o validate.out
./validate.out window_sum_result.txt sum
./validate.out window_median_result.txt median
```

The expected output should be
```
All sliding window [sum|median] are correct.
```

## Folder structure

- manager: redis configuration for the pub-sub manager
- publisher: the publisher impl
- subscriber: two subscribers impl
- output: the output from subscriber and validation script