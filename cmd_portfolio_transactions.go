package main

import (
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"4d63.com/assets/portfolio"
	"github.com/olekukonko/tablewriter"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func cmdPortfolioTransactions(c *ishell.Context, d Data) {
	if len(c.Args) == 0 {
		renderTransactions(os.Stdout, portfolio.Join(d.Portfolios...).Transactions())
	} else {
		transactions := []portfolio.Transaction{}
		for _, arg := range c.Args {
			for _, ah := range portfolio.Join(d.Portfolios...).AggregateHoldings() {
				if strings.ToLower(ah.Stock) == strings.ToLower(arg) {
					transactions = append(transactions, ah.Transactions()...)
				}
			}
		}
		sort.SliceStable(transactions, func(i, j int) bool {
			return transactions[i].Date.Before(transactions[j].Date)
		})
		renderTransactions(os.Stdout, transactions)
	}
}

func cmdPortfolioTransactionsCompleter(args []string, d Data) []string {
	suggestions := []string{}
	for _, ah := range portfolio.Join(d.Portfolios...).AggregateHoldings() {
		stock := strings.ToLower(ah.Stock)
		alreadySelected := false
		for _, arg := range args {
			if stock == strings.ToLower(arg) {
				alreadySelected = true
				break
			}
		}
		if !alreadySelected {
			suggestions = append(suggestions, stock)
		}
	}
	return suggestions
}

func renderTransactions(w io.Writer, transactions []portfolio.Transaction) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Date",
		"Stock",
		"Action",
		"Buy Date",
		"Quantity",
		"Price",
		"Total",
		"Fees",
	})

	for _, t := range transactions {
		buyDate := ""
		if !t.BuyDate.IsZero() {
			buyDate = t.BuyDate.Format("2006-01-02")
		}
		table.Append([]string{
			t.Date.Format("2006-01-02"),
			t.Stock,
			t.Action.String(),
			buyDate,
			strconv.Itoa(t.Quantity),
			t.UnitPrice.StringFixed(2),
			t.TotalPrice().StringFixed(2),
			t.BrokerageFees.StringFixed(2),
		})
	}

	table.Render()
}
