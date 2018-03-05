package exchangerates

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func Read(er ExchangeRates, src, dst string, r io.Reader) error {
	csvReader := csv.NewReader(r)

	i := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading exchange rate %d from csv: %v", i, err)
		}

		for f, field := range record {
			record[f] = strings.TrimSpace(field)
		}

		date, err := time.Parse("02 Jan 2006", record[0])
		if err != nil {
			return fmt.Errorf("reading exchange rate %d from csv: parsing date: %v", i, err)
		}
		rate, err := decimal.NewFromString(record[1])
		if err != nil {
			return fmt.Errorf("reading transaction %d from csv: parsing rate: %v", i, err)
		}

		er.Add(src, dst, date, rate)
	}

	return nil
}
