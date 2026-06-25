// [YOUR_ID_NUMBER] Elad Katz
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	// 1. Parse command-line arguments
	nPtr := flag.Int("n", 0, "total number of orders")
	rPtr := flag.Int("restaurants", 0, "number of restaurants")
	zPtr := flag.Int("zones", 0, "number of delivery zones")
	tokensPtr := flag.String("tokens", "", "comma separated tokens per zone")
	seedAPtr := flag.Int64("seedA", 0, "seed for order generation")
	seedBPtr := flag.Int64("seedB", 0, "seed for processing delays")

	flag.Parse()

	N := *nPtr
	R := *rPtr
	Z := *zPtr
	seedA := *seedAPtr
	seedB := *seedBPtr

	if N == 0 || R == 0 || Z == 0 || *tokensPtr == "" {
		fmt.Println("Missing required arguments")
		os.Exit(1)
	}

	// Parse comma-separated tokens
	tokenStrings := strings.Split(*tokensPtr, ",")
	if len(tokenStrings) != Z {
		fmt.Println("Number of token limits must match number of zones")
		os.Exit(1)
	}

	tokensLimits := make([]int, Z)
	for i, tStr := range tokenStrings {
		tStr = strings.TrimSpace(tStr)
		limit, err := strconv.Atoi(tStr)
		if err != nil {
			fmt.Printf("Invalid token limit: %s\n", tStr)
			os.Exit(1)
		}
		tokensLimits[i] = limit
	}

	// 2. Initialize Core Channels
	events := make(chan Event, N*4) // Buffered logger stream
	done := make(chan bool)

	// Start Logger
	go logger(events, N, done)

	dispatcherChan := make(chan Order, N)
	zoneChannels := make([]chan Order, Z)
	for i := 0; i < Z; i++ {
		zoneChannels[i] = make(chan Order, N)
	}

	// 3. Start Workers (Delivery Zones)
	var zonesWg sync.WaitGroup
	for i := 0; i < Z; i++ {
		zonesWg.Add(1)
		go deliveryZone(i, tokensLimits[i], seedB, zoneChannels[i], events, &zonesWg)
	}

	// 4. Start Producers (Restaurants)
	var restWg sync.WaitGroup
	for i := 0; i < R; i++ {
		restWg.Add(1)
		go restaurant(i, N, R, Z, seedA, dispatcherChan, events, &restWg)
	}

	// 5. Close Dispatcher channel when Producers finish
	go func() {
		restWg.Wait()
		close(dispatcherChan) // Only the creator of the channel closes it
	}()

	// 6. Start Dispatcher
	go runDispatcher(dispatcherChan, zoneChannels, events)

	// 7. Wait for completion
	<-done       // Blocks until the logger prints "DONE total=<N>"
	zonesWg.Wait() // Ensure workers shut down cleanly before exiting main
}