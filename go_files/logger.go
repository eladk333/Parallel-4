// 322587064 - Elad Katz
package main

import "fmt"

func logger(events <-chan Event, N int, done chan<- bool) {
	completedCount := 0
	for e := range events {
		switch e.Kind {
		case "CREATED":
			fmt.Printf("CREATED order=%d restaurant=%d type=%d\n", e.OrderID, e.RestaurantID, e.Zone)
		case "DISPATCHED":
			fmt.Printf("DISPATCHED order=%d zone=%d\n", e.OrderID, e.Zone)
		case "STARTED":
			fmt.Printf("STARTED order=%d zone=%d\n", e.OrderID, e.Zone)
		case "COMPLETED":
			fmt.Printf("COMPLETED order=%d zone=%d\n", e.OrderID, e.Zone)
			completedCount++
			if completedCount == N {
				fmt.Printf("DONE total=%d\n", N)
				done <- true
				return
			}
		}
	}
}