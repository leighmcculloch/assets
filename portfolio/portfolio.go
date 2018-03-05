package portfolio

import (
	"sort"
)

func Join(portfolios ...Portfolio) Portfolio {
	jp := Portfolio{}
	for _, p := range portfolios {
		for _, t := range p.Transactions() {
			jp.Add(t)
		}
	}
	return jp
}

type Portfolio struct {
	Name         string
	transactions []Transaction
	holdings     map[key]Holding
}

type key struct {
	Stock   string
	BuyDate string
}

func (p *Portfolio) Add(t Transaction) {
	if p.transactions == nil {
		p.transactions = []Transaction{}
	}
	p.transactions = append(p.transactions, t)

	if p.holdings == nil {
		p.holdings = map[key]Holding{}
	}

	switch t.Action {
	case ActionBuy:
		k := key{Stock: t.Stock, BuyDate: t.Date.String()}
		h := Holding{
			Stock: t.Stock,
			Buy:   t,
			Sells: []Transaction{},
		}
		p.holdings[k] = h
	case ActionSell:
		k := key{Stock: t.Stock, BuyDate: t.BuyDate.String()}
		h := p.holdings[k]
		h.Sells = append(h.Sells, t)
		p.holdings[k] = h
	}
}

func (p Portfolio) Holdings() []Holding {
	holdings := make([]Holding, 0, len(p.holdings))
	for _, h := range p.holdings {
		holdings = append(holdings, h)
	}
	sort.SliceStable(holdings, func(i, j int) bool {
		return holdings[i].Stock < holdings[j].Stock
	})
	return holdings
}

func (p Portfolio) AggregateHoldings() []AggregateHolding {
	m := map[string]AggregateHolding{}
	for _, h := range p.holdings {
		ah := m[h.Stock]
		ah.Stock = h.Stock
		ah.Add(h)
		m[h.Stock] = ah
	}
	s := make([]AggregateHolding, 0, len(m))
	for _, ah := range m {
		s = append(s, ah)
	}
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Stock < s[j].Stock
	})
	return s
}

func (p Portfolio) Transactions() []Transaction {
	return p.transactions
}
