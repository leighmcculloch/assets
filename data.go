package main

import (
	"4d63.com/assets/exchangerates"
	"4d63.com/assets/portfolio"
)

type Data struct {
	ExchangeRates exchangerates.ExchangeRates
	Portfolios    []portfolio.Portfolio
}

func (d Data) Names() []string {
	names := []string{}
	for _, p := range d.Portfolios {
		names = append(names, p.Name)
	}
	return names
}

func (d Data) Subset(indexes []int) Data {
	indexesMap := map[int]bool{}
	for _, i := range indexes {
		indexesMap[i] = true
	}
	return Data{
		Portfolios: func() []portfolio.Portfolio {
			portfolios := []portfolio.Portfolio{}
			for i, p := range d.Portfolios {
				if indexesMap[i] {
					portfolios = append(portfolios, p)
				}
			}
			return portfolios
		}(),
	}
}
