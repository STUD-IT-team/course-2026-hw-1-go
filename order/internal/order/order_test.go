package order_test

import (
	"testing"

	"order/internal/order"
)

func TestOrderTransitions(t *testing.T) {
	type action struct {
		name string
		call func(*order.Order) bool
	}
	ready := action{"MarkReady", (*order.Order).MarkReady}
	issue := action{"Issue", (*order.Order).Issue}
	cancel := action{"Cancel", (*order.Order).Cancel}

	cases := []struct {
		name  string
		setup []action
		start order.Status
		do    action
		want  order.Status
	}{
		{"pending/ready", nil, order.StatusPending, ready, order.StatusReady},
		{"pending/issue", nil, order.StatusPending, issue, order.StatusPending},
		{"pending/cancel", nil, order.StatusPending, cancel, order.StatusCancelled},
		{"ready/ready", []action{ready}, order.StatusReady, ready, order.StatusReady},
		{"ready/issue", []action{ready}, order.StatusReady, issue, order.StatusIssued},
		{"ready/cancel", []action{ready}, order.StatusReady, cancel, order.StatusCancelled},
		{"issued/ready", []action{ready, issue}, order.StatusIssued, ready, order.StatusIssued},
		{"issued/issue", []action{ready, issue}, order.StatusIssued, issue, order.StatusIssued},
		{"issued/cancel", []action{ready, issue}, order.StatusIssued, cancel, order.StatusIssued},
		{"cancelled/ready", []action{cancel}, order.StatusCancelled, ready, order.StatusCancelled},
		{"cancelled/issue", []action{cancel}, order.StatusCancelled, issue, order.StatusCancelled},
		{"cancelled/cancel", []action{cancel}, order.StatusCancelled, cancel, order.StatusCancelled},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := order.NewOrder("A-42", "Анна", 1500)
			checkOrder(t, o, order.StatusPending)
			for _, step := range tc.setup {
				if !step.call(o) {
					t.Fatalf("%s() during setup = false, want true", step.name)
				}
			}
			checkOrder(t, o, tc.start)

			wantChanged := tc.start != tc.want
			if changed := tc.do.call(o); changed != wantChanged {
				t.Errorf("%s() = %t, want %t", tc.do.name, changed, wantChanged)
			}
			checkOrder(t, o, tc.want)
		})
	}
}

func checkOrder(t *testing.T, o *order.Order, want order.Status) {
	t.Helper()
	if got := o.Status(); got != want {
		t.Errorf("Status() = %v, want %v", got, want)
	}
	if got := o.IsIssued(); got != (want == order.StatusIssued) {
		t.Errorf("IsIssued() = %t in expected state %v", got, want)
	}
	if id, recipient, amount := o.ID(), o.Recipient(), o.Amount(); id != "A-42" || recipient != "Анна" || amount != 1500 {
		t.Errorf("order data = (%q, %q, %d), want (%q, %q, %d)", id, recipient, amount, "A-42", "Анна", 1500)
	}
}

func TestOrdersHaveIndependentState(t *testing.T) {
	first := order.NewOrder("A-42", "Анна", 0)
	second := order.NewOrder("A-42", "Иван", 100000)

	if !first.Cancel() {
		t.Fatal("first order could not be cancelled")
	}
	if second.Status() != order.StatusPending {
		t.Error("cancelling the first order changed the second order")
	}
	if !second.MarkReady() || !second.Issue() {
		t.Fatal("second order could not be issued")
	}
	if first.Status() != order.StatusCancelled {
		t.Error("issuing the second order changed the first order")
	}
	if first.Recipient() != "Анна" || first.Amount() != 0 || second.Recipient() != "Иван" || second.Amount() != 100000 {
		t.Error("orders did not preserve their own recipient and amount")
	}
}
