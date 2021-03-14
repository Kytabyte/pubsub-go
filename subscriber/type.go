package main

// type defines all base types used in subscriber package

import (
	"time"
)

type entry struct {
	Timestamp string // save timestamp as time.RFC3339Nano avoid any precision problem
	Num       int    // the received number
}

// compare function for entries, used for tracking entries in sorted order for median
func (e entry) Compare(other entry) int {
	if e.Num == other.Num {
		t, _ := time.Parse(time.RFC3339Nano, e.Timestamp)
		ot, _ := time.Parse(time.RFC3339Nano, other.Timestamp)
		duration := t.Sub(ot)
		return int(duration)
	}
	return e.Num - other.Num
}

// Subscriber is an interface of all subscribers should implement
type Subscriber interface {
	Receive(int)
}
