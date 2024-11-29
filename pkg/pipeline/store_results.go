package pipeline

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

// StoreResults simulates storing transformed block data with random delays.
func StoreResults(ctx context.Context, item any) error {
	log.Printf("[Store Results] Stored block: %v", item)

	// Simulate random delay
	delay := time.Duration(rand.NormFloat64()*50+200) * time.Millisecond // Mean = 0.2s, Stddev = 0.05s
	if delay < 0 {
		delay = 0 // Ensure delay is non-negative
	}
	time.Sleep(delay)

	// Simulate random failure
	if item.(int)%3 == 0 { // Fail every 3rd block
		return errors.New("store failure: database unavailable")
	}

	return nil
}
