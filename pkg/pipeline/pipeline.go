package pipeline

import (
	"context"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Start initializes the pipeline and manages stages.
func Start(ctx context.Context) {

	ctx, cancel := context.WithCancel(ctx)

	// Setup channels
	feedChan := make(chan any, 10)
	fetchChan := make(chan any, 10)
	transformChan := make(chan any, 10)
	storeChan := make(chan any, 10)

	// WaitGroup for synchronization
	var wg sync.WaitGroup

	// Define pipeline stages
	fetchStage := NewStage(
		"Fetch Blocks",
		FetchBlocks,
		feedChan,
		fetchChan,
	)
	transformStage := NewStage(
		"Transform Blocks",
		TransformBlocks,
		fetchChan,
		transformChan,
	)
	storeStage := NewStage(
		"Store Results",
		StoreResults,
		transformChan,
		storeChan,
	)

	// Start pipeline stages
	go fetchStage.Run(ctx, &wg)
	go transformStage.Run(ctx, &wg)
	go storeStage.Run(ctx, &wg)

	// Simulate continuous block fetching
	blockNumber := 1
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("[Fetch Simulator] Context cancelled, stopping block generation")
				close(feedChan) // Signal that no more blocks will be sent
				return
			default:
				log.Printf("[Fetch Simulator] Sending block %d", blockNumber)
				feedChan <- blockNumber
				blockNumber++
				// Simulate random delay for block generation
				delay := time.Duration(rand.NormFloat64()*1000+4000) * time.Millisecond // Normal(4s, 1s)
				if delay < 0 {
					delay = 0
				}
				time.Sleep(delay)
			}
		}
	}()

	// Listen for errors
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-fetchStage.ErrChan:
				log.Printf("[Fetch Blocks] Error: %v", err)
				handleError(cancel, err)
				os.Exit(1)
			case err := <-transformStage.ErrChan:
				log.Printf("[Transform Blocks] Error: %v", err)
				handleError(cancel, err)
				os.Exit(1)
			case err := <-storeStage.ErrChan:
				log.Printf("[Store Results] Error: %v", err)
				handleError(cancel, err)
				os.Exit(1)
			}
		}
	}()

	// Wait for pipeline stages to complete
	wg.Wait()
}

// handleError cancels the pipeline if a critical error is encountered.
func handleError(cancel context.CancelFunc, err error) {
	log.Printf("Critical error encountered: %v. Cancelling pipeline.", err)
	cancel()
}
