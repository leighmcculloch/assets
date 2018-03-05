package portfolio

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Action string

func (a Action) String() string {
	return strings.Title(string(a))
}

func ParseAction(s string) Action {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return Action(s)
}

const (
	ActionBuy  Action = "buy"
	ActionSell Action = "sell"
)

type Transaction struct {
	Date          time.Time
	Stock         string
	Action        Action
	BuyDate       time.Time
	Quantity      int
	Currency      string
	UnitPrice     decimal.Decimal
	BrokerageFees decimal.Decimal
}

func (t Transaction) TotalPrice() decimal.Decimal {
	return t.UnitPrice.Mul(decimal.New(int64(t.Quantity), 0))
}

func ReadPortfolioTransactions(r io.Reader) ([]Transaction, error) {
	csvReader := csv.NewReader(r)

	transactions := []Transaction{}
	i := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading transaction %d from csv: %v", i, err)
		}

		for j, field := range record {
			record[j] = strings.TrimSpace(field)
		}

		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, fmt.Errorf("reading transaction %d from csv: parsing date: %v", i, err)
		}
		stock := record[1]
		action := ParseAction(record[2])
		buyDate := time.Time{}
		if record[3] != "" {
			buyDate, err = time.Parse("2006-01-02", record[3])
			if err != nil {
				return nil, fmt.Errorf("reading transaction %d from csv: parsing buy date: %v", i, err)
			}
		}
		quantity, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, fmt.Errorf("reading transaction %d from csv: parsing quantity: %v", i, err)
		}
		currency := record[5]
		unitPrice, err := decimal.NewFromString(record[6])
		if err != nil {
			return nil, fmt.Errorf("reading transaction %d from csv: parsing unit price: %v", i, err)
		}
		brokerageFees, err := decimal.NewFromString(record[7])
		if err != nil {
			return nil, fmt.Errorf("reading transaction %d from csv: parsing brokerage fees: %v", i, err)
		}

		transactions = append(
			transactions,
			Transaction{
				Date:          date,
				Stock:         stock,
				Action:        action,
				BuyDate:       buyDate,
				Quantity:      quantity,
				Currency:      currency,
				UnitPrice:     unitPrice,
				BrokerageFees: brokerageFees,
			},
		)
	}

	return transactions, nil
}
