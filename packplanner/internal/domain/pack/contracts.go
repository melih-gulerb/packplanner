package pack

import "context"

// Repository stores the active set of pack sizes used by the application.
type Repository interface {
	List(ctx context.Context) ([]int, error)
	Replace(ctx context.Context, packSizes []int) error
}

// Planner calculates the best shipment plan for a given order quantity.
type Planner interface {
	Calculate(orderQuantity int, packSizes []int) (ShipmentPlan, error)
}
