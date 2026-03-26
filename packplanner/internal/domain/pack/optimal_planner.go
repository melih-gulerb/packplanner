package pack

import (
	"math"
	"sort"
)

type OptimalPlanner struct{}

// NewOptimalPlanner creates the default planner implementation for pack optimization.
func NewOptimalPlanner() OptimalPlanner {
	return OptimalPlanner{}
}

// Calculate applies the business rules in priority order:
// 1. Fulfill the order with the smallest possible total number of items.
// 2. Among those totals, use the lowest number of packs.
func (OptimalPlanner) Calculate(orderQuantity int, packSizes []int) (ShipmentPlan, error) {
	if orderQuantity <= 0 {
		return ShipmentPlan{}, ErrInvalidOrderQuantity
	}

	normalizedPackSizes, err := NormalizePackSizes(packSizes)
	if err != nil {
		return ShipmentPlan{}, err
	}

	smallestPackSize := normalizedPackSizes[0]

	// Any optimal result must be between the order quantity and
	// order quantity + the smallest pack size - 1.
	maxSearchTotal := orderQuantity + smallestPackSize - 1

	bestPackCount, previousTotal, selectedPackSize := initializePlannerState(maxSearchTotal)

	// Totals below the smallest pack size are unreachable, so DP starts from the smallest pack size.
	for total := smallestPackSize; total <= maxSearchTotal; total++ {
		for _, packSize := range normalizedPackSizes {
			if packSize > total {
				break
			}

			if !canReachTotal(total, packSize, bestPackCount) {
				continue
			}

			candidatePackCount := bestPackCount[total-packSize] + 1
			if shouldUseCandidate(candidatePackCount, bestPackCount[total], packSize, selectedPackSize[total]) {
				bestPackCount[total] = candidatePackCount
				previousTotal[total] = total - packSize
				selectedPackSize[total] = packSize
			}
		}
	}

	// Pick the first reachable total at or above the order to minimize over-shipment.
	bestFulfillmentTotal := findBestFulfillmentTotal(orderQuantity, maxSearchTotal, bestPackCount)
	if bestFulfillmentTotal == -1 {
		return ShipmentPlan{}, ErrInvalidPackSizes
	}

	// Rebuild the chosen path from the recorded previous totals.
	packs := buildShipmentPacks(bestFulfillmentTotal, previousTotal, selectedPackSize)

	return ShipmentPlan{
		OrderQuantity: orderQuantity,
		TotalItems:    bestFulfillmentTotal,
		TotalPacks:    bestPackCount[bestFulfillmentTotal],
		Packs:         packs,
	}, nil
}

func initializePlannerState(maxSearchTotal int) ([]int, []int, []int) {
	bestPackCount := make([]int, maxSearchTotal+1)
	previousTotal := make([]int, maxSearchTotal+1)
	selectedPackSize := make([]int, maxSearchTotal+1)

	for total := 1; total <= maxSearchTotal; total++ {
		// MaxInt works as a sentinel for totals that have not been reached yet.
		bestPackCount[total] = math.MaxInt
		previousTotal[total] = -1
	}

	return bestPackCount, previousTotal, selectedPackSize
}

func canReachTotal(total, packSize int, bestPackCount []int) bool {
	return packSize <= total && bestPackCount[total-packSize] != math.MaxInt
}

func shouldUseCandidate(candidatePackCount, currentBestPackCount, candidatePackSize, currentPackSize int) bool {
	if candidatePackCount < currentBestPackCount {
		return true
	}

	// Keep ties deterministic by preferring the larger pack size.
	return candidatePackCount == currentBestPackCount && candidatePackSize > currentPackSize
}

func findBestFulfillmentTotal(orderQuantity, maxSearchTotal int, bestPackCount []int) int {
	for total := orderQuantity; total <= maxSearchTotal; total++ {
		if bestPackCount[total] != math.MaxInt {
			return total
		}
	}

	return -1
}

func buildShipmentPacks(bestFulfillmentTotal int, previousTotal, selectedPackSize []int) []ShipmentPack {
	packCounts := make(map[int]int)
	for total := bestFulfillmentTotal; total > 0; total = previousTotal[total] {
		packCounts[selectedPackSize[total]]++
	}

	packs := make([]ShipmentPack, 0, len(packCounts))
	for packSize, quantity := range packCounts {
		packs = append(packs, ShipmentPack{
			PackSize: packSize,
			Quantity: quantity,
		})
	}

	sort.Slice(packs, func(i, j int) bool {
		return packs[i].PackSize > packs[j].PackSize
	})

	return packs
}
