package main

import (
	"io"
	"os"
	"strconv"

	"4d63.com/assets/portfolio"
	"github.com/olekukonko/tablewriter"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func cmdPortfolioHoldings(c *ishell.Context, d Data) {
	renderHoldings(os.Stdout, portfolio.Join(d.Portfolios...))
}

func renderHoldings(w io.Writer, p portfolio.Portfolio) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Stock",
		"Quantity",
	})

	for _, ah := range p.AggregateHoldings() {
		if ah.Quantity() == 0 {
			continue
		}
		table.Append([]string{
			ah.Stock,
			strconv.Itoa(ah.Quantity()),
		})
	}

	table.Render()
}
