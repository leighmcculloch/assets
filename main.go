package main // import "4d63.com/assets"
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"4d63.com/assets/exchangerates"
	"4d63.com/assets/portfolio"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func main() {
	if err := app(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func app() error {
	shell := ishell.New()
	shell.DeleteCmd("clear")
	shell.DeleteCmd("exit")

	data, err := loadData()
	if err != nil {
		return err
	}

	selectedData := data

	shell.AddCmd(&ishell.Cmd{
		Name: "files",
		Help: "choose files",
		Func: func(c *ishell.Context) {
			selectedData = cmdFiles(c, data)
		},
	})

	portfolioCmd := &ishell.Cmd{
		Name: "portfolio",
		Help: "shares/stock portfolios",
	}
	shell.AddCmd(portfolioCmd)
	portfolioCmd.AddCmd(&ishell.Cmd{
		Name: "transactions",
		Help: "",
		Func: func(c *ishell.Context) {
			cmdPortfolioTransactions(c, selectedData)
		},
		Completer: func(args []string) []string {
			return cmdPortfolioTransactionsCompleter(args, selectedData)
		},
	})
	portfolioCmd.AddCmd(&ishell.Cmd{
		Name: "holdings",
		Help: "access shares/stocks that are held",
		Func: func(c *ishell.Context) {
			cmdPortfolioTransactions(c, selectedData)
		},
	})
	portfolioCmd.AddCmd(&ishell.Cmd{
		Name: "sales",
		Func: func(c *ishell.Context) {
			cmdPortfolioSales(c, selectedData)
		},
	})

	shell.Run()

	return nil
}

func loadData() (Data, error) {
	exchangeRateFiles := []string{}
	portfolioFiles := []string{}

	err := filepath.Walk(".", func(path string, f os.FileInfo, _ error) error {
		if f.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".exchangerate.csv") {
			exchangeRateFiles = append(exchangeRateFiles, path)
		}
		if strings.HasSuffix(path, ".portfolio.csv") {
			portfolioFiles = append(portfolioFiles, path)
		}
		return nil
	})
	if err != nil {
		return Data{}, fmt.Errorf("walking files in current directory: %v", err)
	}

	exchangeRates, err := loadExchangeRates(exchangeRateFiles)
	if err != nil {
		return Data{}, err
	}

	portfolios, err := loadPortfolios(portfolioFiles)
	if err != nil {
		return Data{}, err
	}

	return Data{
		ExchangeRates: exchangeRates,
		Portfolios:    portfolios,
	}, nil
}

func loadExchangeRates(filePaths []string) (exchangerates.ExchangeRates, error) {
	exchangeRates := exchangerates.ExchangeRates{}

	for _, fp := range filePaths {
		var src, dst string
		fmt.Sscanf(fp, "%s-%s.", &src, &dst)
		f, err := os.Open(fp)
		if err != nil {
			return nil, err
		}
		err = exchangerates.Read(exchangeRates, src, dst, f)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %v", fp, err)
		}
	}

	return exchangeRates, nil
}

func loadPortfolios(filePaths []string) ([]portfolio.Portfolio, error) {
	portfolios := []portfolio.Portfolio{}

	for _, fp := range filePaths {
		p, err := loadPortfolio(fp)
		if err != nil {
			return nil, fmt.Errorf("loading portfolio %s: %v", fp, err)
		}
		portfolios = append(portfolios, p)
	}

	return portfolios, nil
}

func loadPortfolio(filePath string) (portfolio.Portfolio, error) {
	transactionsFile, err := os.Open(filePath)
	if err != nil {
		return portfolio.Portfolio{}, fmt.Errorf("opening transactions.csv: %v", err)
	}
	transactions, err := portfolio.ReadPortfolioTransactions(transactionsFile)
	if err != nil {
		return portfolio.Portfolio{}, fmt.Errorf("reading transactions.csv: %v", err)
	}

	p := portfolio.Portfolio{Name: filePath}
	for _, t := range transactions {
		p.Add(t)
	}

	return p, nil
}
