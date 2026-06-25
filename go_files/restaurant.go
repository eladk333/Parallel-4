// 322587064 - Elad Katz
package main

import (
	"math/rand"
	"sync"
)

func restaurant(rid int, N int, R int, Z int, seedA int64, dispatcherChan chan<- Order, events chan<- Event, wg *sync.WaitGroup) {
	defer wg.Done()

	// Deterministic RNG per restaurant
	rng := rand.New(rand.NewSource(seedA + int64(rid)*1000003))

	// Contiguous order IDs
	startID := (rid * N) / R
	endID := ((rid + 1) * N) / R

	for id := startID; id < endID; id++ {
		zone := rng.Intn(Z)
		order := Order{OrderID: id, RestaurantID: rid, FoodType: zone}

		events <- Event{Kind: "CREATED", OrderID: id, RestaurantID: rid, Zone: zone}
		dispatcherChan <- order
	}
}