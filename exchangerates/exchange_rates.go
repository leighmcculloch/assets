package exchangerates

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRates map[exchangeRatesKey]decimal.Decimal

type exchangeRatesKey struct {
	src string
	dst string
}

func (er *ExchangeRates) Add(src, dst string, date time.Time, rate decimal.Decimal) {
	if *er == nil {
		*er = map[exchangeRatesKey]decimal.Decimal{}
	}
	key := exchangeRatesKey{src: src, dst: dst}
	(*er)[key] = rate
}

func (er ExchangeRates) Convert(src, dst string, date time.Time, amount decimal.Decimal) (decimal.Decimal, error) {
	if er == nil {
		return decimal.Zero, fmt.Errorf("no exchange rates defined")
	}

	{
		key := exchangeRatesKey{src: src, dst: dst}
		rate, ok := er[key]
		if ok {
			return amount.Mul(rate), nil
		}
	}

	{
		key := exchangeRatesKey{src: dst, dst: src}
		rate, ok := er[key]
		if ok {
			return amount.Div(rate), nil
		}
	}

	return decimal.Zero, fmt.Errorf("missing exchange rate for %v=>%v on date %v", src, dst, date)
}
