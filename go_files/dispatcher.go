// 322587064 - Elad Katz
package main

func runDispatcher(dispatcherChan <-chan Order, zoneChannels []chan Order, events chan<- Event) {
	// Read from all restaurants (Fan-in) and route to specific zones (Fan-out)
	for order := range dispatcherChan {
		events <- Event{Kind: "DISPATCHED", OrderID: order.OrderID, Zone: order.FoodType}
		zoneChannels[order.FoodType] <- order
	}

	// Close zone channels after dispatching all orders
	for _, zc := range zoneChannels {
		close(zc)
	}
}