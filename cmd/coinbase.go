package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KalebHawkins/crypto-client/coinbase"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// coinbaseCmd represents the coinbase command
var coinbaseCmd = &cobra.Command{
	Use:   "coinbase",
	Short: "interact with the Coinbase API.",
	Long: `Interact with the Coinbase API.

To start working with the Coinbase API you must create an API key and a secret.
You can create an API secret and key by going here: https://www.coinbase.com/settings/api. 
After signing into your account click '+ New API Key' and create your new key.

You will be asked to configure permissions and access to your accounts. Read access for everything 
is sufficient if you do not plan on buying or selling crypto currency using the API. 

Take note of the API key and secret. Now you are ready to use the cypto-client cli. 
To set crypto-client to use your API key and secret export the COINBASE_KEY and COINBASE_SECRET
environment variables. To do this see the examples below.

	[Linux]
	export COINBASE_KEY="API_KEY"
	export COINBASE_SECRET="API_SECRET"

	[Windows (Powershell)]
	$env:COINBASE_KEY = "API_KEY"
	$env:COINBASE_SECRET = "API_SECRET"

Supported operations for Coinbase is depicted in the table below.

	╔═════════════════════════════════════════╤══════════════════╗
	║ Operation                               │ Supported        ║
	╠═════════════════════════════════════════╪══════════════════╣
	║ List high level view of owned assets    │ yes              ║
	╟─────────────────────────────────────────┼──────────────────╢
	║ List transaction data                   │ yes              ║
	╟─────────────────────────────────────────┼──────────────────╢
	║ List account information                │ yes              ║
	╟─────────────────────────────────────────┼──────────────────╢
	║ Buy crypto                              │ work in progress ║
	╟─────────────────────────────────────────┼──────────────────╢
	║ Sell crypto                             │ work in progress ║
	╟─────────────────────────────────────────┼──────────────────╢
	║ Set profile information                 │ work in progress ║
	╚═════════════════════════════════════════╧══════════════════╝
`,

	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		if listTransactions {
			getCoinbaseTransactions()
		}

		if listAccounts {
			getCoinbaseAccounts()
		}

		if !listAccounts && !listTransactions {
			getCoinbaseOverview()
		}

		fmt.Println()
		fmt.Println("Elapsed Run Time:", time.Since(start))
	},
}

var listTransactions bool
var listAccounts bool

func init() {
	rootCmd.AddCommand(coinbaseCmd)
	coinbaseCmd.Flags().BoolVarP(&listTransactions, "list-transactions", "t", false, "list all your accounts transactions")
	coinbaseCmd.Flags().BoolVarP(&listAccounts, "list-accounts", "a", false, "list all your accounts")
}

// getCoinbaseOverview will output a wholistic overview of your Coinbase account and assets.
// This is the default when running `crypto-client coinbase` without additional flags.
func getCoinbaseOverview() {
	c := coinbase.APIKeyClient()
	user, err := c.GetUserProfile()
	errHandler(err)
	fmt.Println(user)

	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()

	tbl := table.New("Wallet", "Balance", "Currency", "Spot Price Per Unit",
		"Buy Price Per Unit", "Sell Price Per Unit", "Total Sell Out Price", "Invested",
		"Inflation Rewards", "Total Return")
	tbl.WithHeaderFormatter(headerFmt)

	account, err := c.GetAccount()
	errHandler(err)

	var totalSellOutAmount float64
	var totalReturnAmount float64

	for _, act := range account.Data {
		amt, err := strconv.ParseFloat(act.Balance.Amount, 64)
		errHandler(err)

		if amt > 0 {

			currencyPair := fmt.Sprintf("%s-%s", act.Balance.Currency, user.Data.NativeCurrency)

			spotPrice, err := c.GetPrice(currencyPair, coinbase.Spot)
			errHandler(err)
			spotAmt, err := strconv.ParseFloat(spotPrice.Data.Amount, 64)
			errHandler(err)
			buyPrice, err := c.GetPrice(currencyPair, coinbase.Buy)
			errHandler(err)
			bpAmt, err := strconv.ParseFloat(buyPrice.Data.Amount, 64)
			errHandler(err)

			sellPrice, err := c.GetPrice(currencyPair, coinbase.Sell)
			errHandler(err)
			sellAmt, err := strconv.ParseFloat(sellPrice.Data.Amount, 64)
			errHandler(err)

			var invested float64
			var inflationRewards float64

			transactions, err := c.GetTransactionHistory(act.ID)
			errHandler(err)

			for _, tr := range transactions.Data {
				trNcAmt, err := strconv.ParseFloat(tr.NativeAmount.Amount, 64)
				errHandler(err)
				trAmt, err := strconv.ParseFloat(tr.Amount.Amount, 64)
				errHandler(err)

				switch tr.Type {
				case coinbase.Buy:
					invested += trNcAmt
				case coinbase.InflationReward:
					inflationRewards += trAmt
				}

			}

			sellOutAmount := amt * sellAmt
			returnAmount := sellOutAmount - invested

			tbl.AddRow(act.Name, fmt.Sprintf("%f", amt), act.Balance.Currency,
				fmt.Sprintf("%.2f %s", spotAmt, spotPrice.Data.Currency),
				fmt.Sprintf("%.2f %s", bpAmt, buyPrice.Data.Currency),
				fmt.Sprintf("%.2f %s", sellAmt, sellPrice.Data.Currency),
				fmt.Sprintf("%.2f %s", sellOutAmount, sellPrice.Data.Currency),
				fmt.Sprintf("%.2f %s", invested, user.Data.NativeCurrency),
				fmt.Sprintf("%f %s", inflationRewards, act.Balance.Currency),
				fmt.Sprintf("%.2f %s", returnAmount, user.Data.NativeCurrency))

			totalSellOutAmount += amt * sellAmt
			totalReturnAmount += returnAmount

		}
	}

	tbl.Print()

	fmt.Printf("Total Sell Out Amount: %.2f %s\n", totalSellOutAmount, user.Data.NativeCurrency)
	fmt.Printf("Total Return Amount: %.2f %s\n", totalReturnAmount, user.Data.NativeCurrency)
}

// getCoinbaseTransactions will list all past transactions the currency and a summary.
func getCoinbaseTransactions() {
	table.DefaultHeaderFormatter = func(s string, i ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(s, i...))
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	tbl := table.New("Transaction Type", "Crypto", "Amount", "Date", "Payment Method", "Summary").WithHeaderFormatter(headerFmt)

	c := coinbase.APIKeyClient()

	accounts, err := c.GetAccount()
	errHandler(err)

	var wg sync.WaitGroup
	for _, a := range accounts.Data {
		wg.Add(1)
		go func(accountID string) {
			defer wg.Done()
			tr, err := c.GetTransactionHistory(accountID)
			errHandler(err)

			for _, t := range tr.Data {
				tAmt, err := strconv.ParseFloat(t.Amount.Amount, 64)
				errHandler(err)

				tbl.AddRow(t.Type, t.Amount.Currency, tAmt, t.CreatedAt, t.Details.PaymentMethodName, t.Details.Header)
			}
		}(a.ID)
	}
	wg.Wait()

	tbl.Print()
}

// getCoinbaseAccounts will list all your coinbase accounts that contain assets.
func getCoinbaseAccounts() {

	table.DefaultHeaderFormatter = func(s string, i ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(s, i...))
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	tbl := table.New("Wallet", "Balance", "Native").WithHeaderFormatter(headerFmt)

	c := coinbase.APIKeyClient()
	user, err := c.GetUserProfile()
	errHandler(err)

	acts, err := c.GetAccount()
	errHandler(err)

	var wg sync.WaitGroup
	wg.Add(len(acts.Data))

	for _, a := range acts.Data {
		amt, err := strconv.ParseFloat(a.Balance.Amount, 64)
		errHandler(err)
		if amt > 0 {
			currencyPair := fmt.Sprintf("%s-%s", a.Balance.Currency, user.Data.NativeCurrency)
			spotPrice, err := c.GetPrice(currencyPair, coinbase.Spot)
			errHandler(err)
			sAmt, err := strconv.ParseFloat(spotPrice.Data.Amount, 64)
			errHandler(err)

			tbl.AddRow(a.Name, a.Balance.Amount, fmt.Sprintf("%.2f %s", sAmt*amt, user.Data.NativeCurrency))
		}
	}

	tbl.Print()
}

// errHandler is a short hand error handler.
func errHandler(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
}
