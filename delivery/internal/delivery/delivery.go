package delivery

type Tariff interface {
	Cost(weightGrams int) int
}

type Pickup struct{}

func (Pickup) Cost(weightGrams int) int {
	return 0 // TODO: Реализуйте тариф самовывоза.
}

type Courier struct{}

func (Courier) Cost(weightGrams int) int {
	return 0 // TODO: Рассчитайте стоимость по весу.
}

func TotalDelivery(weights []int, tariff Tariff) int {
	// TODO: Рассчитайте общую стоимость через переданный тариф.
	return 0
}
