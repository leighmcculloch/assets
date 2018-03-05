package main

import (
	"io"
	"os"
	"strconv"

	"4d63.com/assets/portfolio"
	"github.com/olekukonko/tablewriter"
	"github.com/shopspring/decimal"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func cmdPortfolioSales(c *ishell.Context, d Data) {
	renderSales(os.Stdout, portfolio.Join(d.Portfolios...))
}

func renderSales(w io.Writer, p portfolio.Portfolio) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Stock",
		"Date",
		"Quantity",
		"Cost",
		"Sell",
		"Gain",
		"Loss",
	})

	for _, h := range p.Holdings() {
		if len(h.Sells) == 0 {
			continue
		}

		for _, s := range h.Sells {
			cost := h.Buy.UnitPrice.Mul(decimal.New(int64(s.Quantity), 0))
			sell := s.UnitPrice.Mul(decimal.New(int64(s.Quantity), 0))
			gain := sell.Sub(cost)
			loss := decimal.Zero
			if gain.LessThan(decimal.Zero) {
				loss = gain.Abs()
				gain = decimal.Zero
			}
			table.Append([]string{
				s.Stock,
				s.Date.Format("2006-01-02"),
				strconv.Itoa(s.Quantity),
				cost.StringFixed(2),
				sell.StringFixed(2),
				func() string {
					if gain.Equal(decimal.Zero) {
						return ""
					} else {
						return gain.StringFixed(2)
					}
				}(),
				func() string {
					if loss.Equal(decimal.Zero) {
						return ""
					} else {
						return loss.StringFixed(2)
					}
				}(),
			})
		}
	}

	table.Render()
}
