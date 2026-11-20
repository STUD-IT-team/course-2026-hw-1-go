package order_test

import (
	"slices"
	"testing"

	"order/internal/order"
)

func TestIssueBatch(t *testing.T) {
	pending := order.NewOrder("pending", "Анна", 10)
	cancelled := order.NewOrder("cancelled", "Анна", 10)
	cancelled.Cancel()
	issued := newReadyOrder(t, "issued", 10)
	issued.Issue()

	first := newReadyOrder(t, "A", 700)
	expensive := newReadyOrder(t, "B", 400)
	fits := newReadyOrder(t, "C", 300)
	free := newReadyOrder(t, "D", 0)
	sameID := newReadyOrder(t, "D", 0)
	orders := []*order.Order{pending, first, expensive, cancelled, fits, issued, free, free, sameID}

	ids, total := order.IssueBatch(orders, 1000)
	if want := []string{"A", "C", "D", "D"}; !slices.Equal(ids, want) || total != 1000 {
		t.Errorf("IssueBatch() = (%v, %d), want (%v, 1000)", ids, total, want)
	}

	wantStatuses := []order.Status{
		order.StatusPending, order.StatusIssued, order.StatusReady,
		order.StatusCancelled, order.StatusIssued, order.StatusIssued,
		order.StatusIssued, order.StatusIssued, order.StatusIssued,
	}
	for i, o := range orders {
		if got := o.Status(); got != wantStatuses[i] {
			t.Errorf("orders[%d].Status() = %v, want %v", i, got, wantStatuses[i])
		}
	}

	ids, total = order.IssueBatch(orders, 400)
	if !slices.Equal(ids, []string{"B"}) || total != 400 || !expensive.IsIssued() {
		t.Errorf("second IssueBatch() = (%v, %d), B issued = %t, want ([B], 400), B issued = true",
			ids, total, expensive.IsIssued())
	}
}

func TestIssueBatchWithoutEligibleOrders(t *testing.T) {
	paid := newReadyOrder(t, "A", 1)
	for _, orders := range [][]*order.Order{nil, {paid}} {
		ids, total := order.IssueBatch(orders, 0)
		if len(ids) != 0 || total != 0 {
			t.Errorf("IssueBatch() = (%v, %d), want ([], 0)", ids, total)
		}
	}
	if paid.Status() != order.StatusReady {
		t.Error("order above the limit should remain ready")
	}
}

func newReadyOrder(t *testing.T, id string, amount int) *order.Order {
	t.Helper()
	o := order.NewOrder(id, "Анна", amount)
	if !o.MarkReady() {
		t.Fatalf("could not mark order %q ready", id)
	}
	return o
}
