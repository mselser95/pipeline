package pipeline

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

// TransformBlocks simulates transforming block data with random delays.
func TransformBlocks(ctx context.Context, item any) error {
	log.Printf("[Transform Blocks] Transformed block: %v", item)

	// Simulate random delay
	delay := time.Duration(rand.NormFloat64()*500+2000) * time.Millisecond // Mean = 2s, Stddev = 0.5s
	if delay < 0 {
		delay = 0 // Ensure delay is non-negative
	}
	time.Sleep(delay)

	// Simulate random failure
	if item.(int)%5 == 0 { // Fail every 5th block
		return errors.New("transform failure: data inconsistency")
	}

	return nil
}
