// [YOUR_ID_NUMBER] Elad Katz
package main

import (
	"math/rand"
	"sync"
	"time"
)

func deliveryZone(zoneID int, tokensLimit int, seedB int64, zoneChan <-chan Order, events chan<- Event, wg *sync.WaitGroup) {
	defer wg.Done()

	// Deterministic RNG per zone
	rng := rand.New(rand.NewSource(seedB + int64(zoneID)*2000003))

	// Token limit channel
	tokens := make(chan struct{}, tokensLimit)

	for order := range zoneChan {
		tokens <- struct{}{} // Acquire token (blocks if full)

		events <- Event{Kind: "STARTED", OrderID: order.OrderID, Zone: zoneID}

		// Deterministic processing delay
		delay := time.Duration(rng.Intn(21)) * time.Millisecond
		time.Sleep(delay)

		events <- Event{Kind: "COMPLETED", OrderID: order.OrderID, Zone: zoneID}

		<-tokens // Return token
	}
}