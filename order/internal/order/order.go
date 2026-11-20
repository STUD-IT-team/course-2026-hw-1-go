package order

type Status int

const (
	StatusPending Status = iota
	StatusReady
	StatusIssued
	StatusCancelled
)

type Order struct {
	// TODO: Добавьте закрытые поля для данных и состояния заказа.
}

func NewOrder(id, recipient string, amount int) *Order {
	// TODO: Сохраните данные и установите начальное состояние.
	return &Order{}
}

func (o *Order) ID() string {
	return "" // TODO
}

func (o *Order) Recipient() string {
	return "" // TODO
}

func (o *Order) Amount() int {
	return 0 // TODO
}

func (o *Order) Status() Status {
	return StatusPending // TODO
}

func (o *Order) IsIssued() bool {
	return false // TODO
}

func (o *Order) MarkReady() bool {
	return false // TODO
}

func (o *Order) Issue() bool {
	return false // TODO
}

func (o *Order) Cancel() bool {
	return false // TODO
}
