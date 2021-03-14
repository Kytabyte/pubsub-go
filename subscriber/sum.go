package main

import (
	"container/list"
	"fmt"
	"io"
	"time"
)

// SumSubscriber is a type of subscriber which takes an integer
// and print the sum of numbers that received
type SumSubscriber struct {
	queue  *list.List
	sum    int
	writer io.Writer
}

// NewSumSubscriber creates a new SumSubsriber
func NewSumSubscriber(writer io.Writer) *SumSubscriber {
	return &SumSubscriber{
		queue:  new(list.List),
		sum:    0,
		writer: writer,
	}
}

// Receive an integer and print the sum of numbers received in last 5 seconds
func (s *SumSubscriber) Receive(num int) {
	q := s.queue
	t := time.Now()

	// remove outdated entry
	for q.Len() > 0 {
		pt, _ := time.Parse(time.RFC3339Nano, q.Front().Value.(entry).Timestamp)
		if t.Sub(pt) <= 5*time.Second {
			break
		}
		e := q.Remove(q.Front()).(entry)
		s.sum -= e.Num
	}

	// push new entry
	e := entry{
		Timestamp: t.Format(time.RFC3339Nano),
		Num:       num,
	}
	q.PushBack(e)
	s.sum += num
	fmt.Fprintf(s.writer, "Received %d at %v, Sum %d.\n", num, t.Format(time.RFC3339Nano), s.sum)
}
