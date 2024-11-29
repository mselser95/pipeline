package pipeline

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

// FetchBlocks simulates fetching block data with random delays.
func FetchBlocks(ctx context.Context, item any) error {
	blockNumber := item.(int)
	log.Printf("[Fetch Blocks] Fetched block %d", blockNumber)

	// Simulate random delay
	delay := time.Duration(rand.NormFloat64()*1000+4000) * time.Millisecond // Mean = 4s, Stddev = 1s
	if delay < 0 {
		delay = 0 // Ensure delay is non-negative
	}
	time.Sleep(delay)

	// Simulate random failure
	if blockNumber%7 == 0 { // Fail every 7th block
		return errors.New("fetch failure: unable to fetch block")
	}

	return nil
}
