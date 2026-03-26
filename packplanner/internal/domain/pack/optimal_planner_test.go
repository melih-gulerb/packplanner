package pack

import "testing"

func TestOptimalPlannerCalculate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		orderQuantity int
		packSizes     []int
		wantItems     int
		wantPacks     int
		wantBreakdown map[int]int
	}{
		{
			name:          "single smallest pack when below minimum",
			orderQuantity: 1,
			packSizes:     []int{250, 500, 1000, 2000, 5000},
			wantItems:     250,
			wantPacks:     1,
			wantBreakdown: map[int]int{250: 1},
		},
		{
			name:          "single pack size still fulfills below minimum order",
			orderQuantity: 1,
			packSizes:     []int{250},
			wantItems:     250,
			wantPacks:     1,
			wantBreakdown: map[int]int{250: 1},
		},
		{
			name:          "fewer packs chosen for same total",
			orderQuantity: 251,
			packSizes:     []int{250, 500, 1000, 2000, 5000},
			wantItems:     500,
			wantPacks:     1,
			wantBreakdown: map[int]int{500: 1},
		},
		{
			name:          "lower total beats fewer packs",
			orderQuantity: 501,
			packSizes:     []int{250, 500, 1000, 2000, 5000},
			wantItems:     750,
			wantPacks:     2,
			wantBreakdown: map[int]int{500: 1, 250: 1},
		},
		{
			name:          "large order example",
			orderQuantity: 12001,
			packSizes:     []int{250, 500, 1000, 2000, 5000},
			wantItems:     12250,
			wantPacks:     4,
			wantBreakdown: map[int]int{5000: 2, 2000: 1, 250: 1},
		},
		{
			name:          "configurable sizes can beat greedy selection",
			orderQuantity: 1000,
			packSizes:     []int{250, 500, 700},
			wantItems:     1000,
			wantPacks:     2,
			wantBreakdown: map[int]int{500: 2},
		},
	}

	planner := NewOptimalPlanner()

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			plan, err := planner.Calculate(testCase.orderQuantity, testCase.packSizes)
			if err != nil {
				t.Fatalf("Calculate() error = %v", err)
			}

			if plan.TotalItems != testCase.wantItems {
				t.Fatalf("TotalItems = %d, want %d", plan.TotalItems, testCase.wantItems)
			}

			if plan.TotalPacks != testCase.wantPacks {
				t.Fatalf("TotalPacks = %d, want %d", plan.TotalPacks, testCase.wantPacks)
			}

			if len(plan.Packs) != len(testCase.wantBreakdown) {
				t.Fatalf("len(Packs) = %d, want %d", len(plan.Packs), len(testCase.wantBreakdown))
			}

			for _, shipmentPack := range plan.Packs {
				wantQuantity, exists := testCase.wantBreakdown[shipmentPack.PackSize]
				if !exists {
					t.Fatalf("unexpected pack size in result: %d", shipmentPack.PackSize)
				}

				if shipmentPack.Quantity != wantQuantity {
					t.Fatalf("pack size %d quantity = %d, want %d", shipmentPack.PackSize, shipmentPack.Quantity, wantQuantity)
				}
			}
		})
	}
}

func TestNormalizePackSizes(t *testing.T) {
	t.Parallel()

	got, err := NormalizePackSizes([]int{500, 250, 500, 1000})
	if err != nil {
		t.Fatalf("NormalizePackSizes() error = %v", err)
	}

	want := []int{250, 500, 1000}
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}

	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("got[%d] = %d, want %d", index, got[index], want[index])
		}
	}
}
