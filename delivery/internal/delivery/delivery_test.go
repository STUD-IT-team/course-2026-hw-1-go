package delivery_test

import (
	"slices"
	"testing"

	"delivery/internal/delivery"
)

func TestPickupCost(t *testing.T) {
	var tariff delivery.Tariff = delivery.Pickup{}
	for _, weight := range []int{1, 1000, 1001, 30000} {
		if got := tariff.Cost(weight); got != 0 {
			t.Errorf("Pickup.Cost(%d) = %d, want 0", weight, got)
		}
	}
}

func TestCourierCost(t *testing.T) {
	var tariff delivery.Tariff = delivery.Courier{}
	cases := []struct{ weight, want int }{
		{1, 150},
		{999, 150},
		{1000, 150},
		{1001, 200},
		{2000, 200},
		{2001, 250},
		{29001, 1600},
		{30000, 1600},
	}
	for _, tc := range cases {
		if got := tariff.Cost(tc.weight); got != tc.want {
			t.Errorf("Courier.Cost(%d) = %d, want %d", tc.weight, got, tc.want)
		}
	}
}

func TestTotalDeliveryWithBuiltInTariffs(t *testing.T) {
	weights := []int{1000, 1001}
	if got := delivery.TotalDelivery(weights, delivery.Pickup{}); got != 0 {
		t.Errorf("TotalDelivery with Pickup = %d, want 0", got)
	}
	if got := delivery.TotalDelivery(weights, delivery.Courier{}); got != 350 {
		t.Errorf("TotalDelivery with Courier = %d, want 350", got)
	}
}

type perGramTariff struct {
	calls []int
}

func (tariff *perGramTariff) Cost(weightGrams int) int {
	tariff.calls = append(tariff.calls, weightGrams)
	return 2 * weightGrams
}

func TestTotalDeliveryWithCustomTariff(t *testing.T) {
	tariff := &perGramTariff{}
	weights := []int{30, 10, 30}
	if got := delivery.TotalDelivery(weights, tariff); got != 140 {
		t.Errorf("TotalDelivery with custom tariff = %d, want 140", got)
	}
	if want := []int{30, 10, 30}; !slices.Equal(tariff.calls, want) {
		t.Errorf("Cost calls = %v, want %v", tariff.calls, want)
	}
}

func TestTotalDeliveryWithoutOrders(t *testing.T) {
	tariff := &perGramTariff{}
	if got := delivery.TotalDelivery(nil, tariff); got != 0 {
		t.Errorf("TotalDelivery(nil, tariff) = %d, want 0", got)
	}
	if len(tariff.calls) != 0 {
		t.Errorf("Cost called with weights %v for an empty batch", tariff.calls)
	}
}
