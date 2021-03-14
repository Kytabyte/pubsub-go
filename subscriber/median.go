package main

import (
	"container/list"
	"fmt"
	"io"
	"sort"
	"time"
)

// MedianSubscriber is a type of subscriber which takes an integer
// and print the median of numbers that received within last 5 seconds
type MedianSubscriber struct {
	queue   *list.List
	entries []entry // sorted entries for getting median
	writer  io.Writer
}

// NewMedianSubscriber creates a new SumSubsriber
func NewMedianSubscriber(writer io.Writer) *MedianSubscriber {
	return &MedianSubscriber{
		queue:   new(list.List),
		entries: make([]entry, 0),
		writer:  writer,
	}
}

// Receive an integer and print the median of numbers received in last 5 seconds
func (s *MedianSubscriber) Receive(num int) {
	q := s.queue
	t := time.Now()

	// remove outdated entry
	for q.Len() > 0 {
		pt, _ := time.Parse(time.RFC3339Nano, q.Front().Value.(entry).Timestamp)
		if t.Sub(pt) <= 5*time.Second {
			break
		}
		e := q.Remove(q.Front()).(entry)
		s.removeEntry(e)
	}

	e := entry{
		Timestamp: t.Format(time.RFC3339Nano),
		Num:       num,
	}
	q.PushBack(e)
	s.addEntry(e)
	fmt.Fprintf(s.writer, "Received %d at %v, Median %.1f.\n", num, t.Format(time.RFC3339Nano), s.getMedian())
}

// The best way to handle sliding window median is to use two sorted maps (one increasing and the other one decreasing)
// (like map in c++ or TreeMap in Java).
// But Go doesn't have one of these containers. So we have to maintain a sorted slice using O(n) insert and delete
func (s *MedianSubscriber) addEntry(e entry) {
	n := len(s.entries)
	pos := sort.Search(n, func(i int) bool {
		return s.entries[i].Compare(e) >= 0
	}) // insert pos of entry
	if pos == n {
		s.entries = append(s.entries, e)
	} else {
		s.entries = append(s.entries[:pos+1], s.entries[pos:]...)
		s.entries[pos] = e
	}
}

func (s *MedianSubscriber) removeEntry(e entry) {
	for pos := range s.entries {
		if s.entries[pos].Compare(e) == 0 {
			s.entries = append(s.entries[:pos], s.entries[pos+1:]...)
			return
		}
	}
}

func (s *MedianSubscriber) getMedian() float64 {
	n := len(s.entries)
	if n%2 == 1 {
		return float64(s.entries[n/2].Num)
	}
	return float64(s.entries[n/2-1].Num+s.entries[n/2].Num) / 2
}
