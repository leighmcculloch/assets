package portfolio

type Holding struct {
	Stock string
	Buy   Transaction
	Sells []Transaction
}

func (h Holding) Quantity() int {
	q := h.Buy.Quantity
	for _, s := range h.Sells {
		q -= s.Quantity
	}
	return q
}

func (h Holding) Transactions() []Transaction {
	return append([]Transaction{h.Buy}, h.Sells...)
}

type AggregateHolding struct {
	Stock    string
	holdings []Holding
}

func (ah *AggregateHolding) Add(h Holding) {
	if ah.holdings == nil {
		ah.holdings = []Holding{}
	}
	ah.holdings = append(ah.holdings, h)
}

func (ah AggregateHolding) Quantity() int {
	q := 0
	for _, h := range ah.holdings {
		q += h.Buy.Quantity
		for _, s := range h.Sells {
			q -= s.Quantity
		}
	}
	return q
}

func (ah AggregateHolding) Transactions() []Transaction {
	transactions := []Transaction{}
	for _, h := range ah.holdings {
		transactions = append(transactions, h.Transactions()...)
	}
	return transactions
}
