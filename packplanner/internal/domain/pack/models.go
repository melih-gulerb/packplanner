package pack

// ShipmentPlan is the domain result returned after applying the packing rules.
type ShipmentPlan struct {
	OrderQuantity int
	TotalItems    int
	TotalPacks    int
	Packs         []ShipmentPack
}

// ShipmentPack represents how many packs of a specific size are included in the plan.
type ShipmentPack struct {
	PackSize int
	Quantity int
}
