// 322587064 - Elad Katz
package main

type Order struct {
	OrderID      int // unique in [0..N-1]
	RestaurantID int // in [0..R-1]
	FoodType     int // zone index in [0..Z-1]
}

type Event struct {
	Kind         string // "CREATED" | "DISPATCHED" | "STARTED" | "COMPLETED" | "DONE"
	OrderID      int
	RestaurantID int
	Zone         int
}