/*
Package coinbase is used to query the Coinbase API for information on a user's profile, accounts,
transactions, and exchange rates.
*/
package coinbase

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rodaine/table"
)

// APIKeyClient sets the API key and API secret for Coinbase authentication.
// to use your API Key and API secret set your environment variables.
//  export COINBASE_API="api_key"
//  export COINBASE_SECRET="api_secret"
func APIKeyClient() CoinbaseClient {
	cbAPIKey = os.Getenv("COINBASE_KEY")
	cbAPISecret = os.Getenv("COINBASE_SECRET")

	return CoinbaseClient{}
}

// ─── COINBASE METHODS ───────────────────────────────────────────────────────────

// GetUserProfile upon a successful API request returns a user's profile information. An error is returned
// if creating or sending the request failed.
func (c CoinbaseClient) GetUserProfile() (User, error) {

	body, err := createRequest("user")

	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(body, &user)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetAccount upon a successful API request returns coinbase account information. An error is returned
// if creating or sending the request failed.
func (c CoinbaseClient) GetAccount() (Account, error) {

	body, err := createRequest("accounts")

	if err != nil {
		return Account{}, err
	}

	var account Account
	err = json.Unmarshal(body, &account)

	if err != nil {
		return Account{}, err
	}

	return account, nil
}

// GetExchangeRate() upon a successful API request returns coinbase exchange rate information. An error is returned
// if creating or sending the request failed.
func (c CoinbaseClient) GetExchangeRate() (ExchangeRate, error) {
	body, err := createRequest("exchange-rates")

	if err != nil {
		return nil, err
	}

	var exchangeRate ExchangeRate
	err = json.Unmarshal(body, &exchangeRate)

	if err != nil {
		return nil, err
	}

	return exchangeRate, nil
}

// GetPrice() upon a successful API request returns coinbase price information. An error is returned
// if creating or sending the request failed.
// The `currencyPair` parameter is the currency in which you want to get the
// price for. For example passing "BTC-USD" will return the amount of USD it would take to purchase 1 unit if
// BTC.
// The `priceType` parameter is used to determine what price information to lookup. There are three valid types:
//
//   - buy
//   - sell
//   - spot
//
// These string values are mapped using the constant values `coinbase.Buy`, `coinbase.Sell`, and `coinbase.Spot` defined in the `types.go` file.
func (c CoinbaseClient) GetPrice(currencyPair string, priceType string) (Price, error) {
	body, err := createRequest(fmt.Sprintf("prices/%s/%s", currencyPair, priceType))

	if err != nil {
		return Price{}, nil
	}

	var sp Price
	err = json.Unmarshal(body, &sp)

	if err != nil {
		return Price{}, nil
	}
	return sp, nil
}

// GetPriceByDate() upon a successful API request returns coinbase price information. An error is returned
// if creating or sending the request failed.
// The `currencyPair` parameter is the currency in which you want to get the
// price for. For example passing "BTC-USD" will return the amount of USD it would take to purchase 1 unit if
// BTC.
// The `year` is a time object formatted as YYYY-MM-DD.
func (c CoinbaseClient) GetPriceByDate(currencyPair string, year time.Time) (Price, error) {

	body, err := createRequest(fmt.Sprintf("prices/%s/spot?date=%s", currencyPair, year.Format("2006-01-02")))

	if err != nil {
		return Price{}, err
	}

	var p Price
	err = json.Unmarshal(body, &p)

	if err != nil {
		return Price{}, err
	}

	return p, nil
}

// GetTransactionHistory upon a successful API request returns coinbase transaction information. An error is returned
// if creating or sending the request failed. The `accountID` parameter is the account ID in which you want to get the
// transactions for.
func (c CoinbaseClient) GetTransactionHistory(accountId string) (Transaction, error) {
	body, err := createRequest(fmt.Sprintf("accounts/%v/transactions", accountId))

	if err != nil {
		return Transaction{}, err
	}

	var t Transaction
	err = json.Unmarshal(body, &t)

	if err != nil {
		return Transaction{}, nil
	}

	return t, nil
}

//
// ────────────────────────────────────────────────────────── COIBASE METHODS ─────
//

// ─── STRINGER METODS ────────────────────────────────────────────────────────────

// User.String() is a stringer function for a coinbase User object.
// It only displays a subset of information about the User profile.
func (u User) String() string {
	return fmt.Sprintf("Name: %v\nCountry: %v\nState: %v\nTimezone: %v\nNative Curreny: %v\nBitcoinUnit: %v\nAccount Created: %v\n",
		u.Data.Name, u.Data.Country.Code, u.Data.State, u.Data.TimeZone, u.Data.NativeCurrency, u.Data.BitcoinUnit, u.Data.CreatedAt.Local().Format("01-02-2006 15:04"))
}

// Account.String() is a stringer function for a coinbase Account object.
// It only displays a subset of information about accounts.
func (a Account) String() string {
	var buf bytes.Buffer
	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}

	tbl := table.New("Wallet Name", "Balance", "Currency").WithWriter(&buf)

	for _, act := range a.Data {
		amount, err := strconv.ParseFloat(act.Balance.Amount, 64)
		if err != nil {
			panic(err)
		}

		if amount > 0 {
			tbl.AddRow(act.Name, fmt.Sprintf("%f", amount), act.Balance.Currency)
		}
	}
	tbl.Print()

	return buf.String()
}

// ExchangeRate.String() is a stringer function for a coinbase ExchangeRate object.
// It only displays a subset of information about accounts.
func (e ExchangeRate) String() string {
	data := e["data"].(map[string]interface{})
	currency := data["currency"].(string)
	rates := data["rates"].(map[string]interface{})

	var buf bytes.Buffer
	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}
	tbl := table.New("Currency", "Crypto", "Rate").WithWriter(&buf)

	for k, v := range rates {
		tbl.AddRow(currency, k, v.(string))
	}
	tbl.Print()
	return buf.String()
}

// SpotPrice.String() is a stringer function for a coinbase SpotPrice object.
func (p Price) String() string {
	amt, _ := strconv.ParseFloat(p.Data.Amount, 64)
	return fmt.Sprintf("%s: %.2f %s", p.Data.Base, amt, p.Data.Currency)
}

// Transaction.String() is a stringer function for a coinbase Transaction object.
func (tr Transaction) String() string {

	var buf bytes.Buffer
	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}

	tbl := table.New("Transaction Type", "Crypto", "Amount", "Native Curreny", "Amount", "Date", "Payment Method", "Summary").WithWriter(&buf)

	for _, t := range tr.Data {
		cAmt, _ := strconv.ParseFloat(t.Amount.Amount, 64)
		ncAmt, _ := strconv.ParseFloat(t.NativeAmount.Amount, 64)

		tbl.AddRow(t.Type, t.Amount.Currency, cAmt, t.NativeAmount.Currency, ncAmt, t.CreatedAt.Format("2006-01-02 15:04"), t.Details.PaymentMethodName, t.Details.Header)
	}
	tbl.Print()

	return buf.String()
}

//
// ───────────────────────────────────────────────────────── STRINGER METHODS ─────
//

// ─── HELPER FUNCTIONS ───────────────────────────────────────────────────────────

// createSignature returns the sha value for the CB-ACCESS-SIGN header that Coinbase requires for its API calls.
func createSignature(r *http.Request) string {
	timestamp := time.Now().Unix()
	h := hmac.New(sha256.New, []byte(cbAPISecret))
	h.Write([]byte(fmt.Sprintf("%v%v%v", timestamp, r.Method, r.URL.Path)))

	return hex.EncodeToString(h.Sum(nil))
}

// appendHeaders appends the Coinbase required API Headers
func appendHeaders(r *http.Request, sig string) {
	r.Header.Add("CB-ACCESS-KEY", cbAPIKey)
	r.Header.Add("CB-ACCESS-SIGN", sig)
	r.Header.Add("CB-ACCESS-TIMESTAMP", fmt.Sprintf("%v", time.Now().Unix()))
	r.Header.Add("CB-VERSION", cbAPIVersion)
	r.Header.Add("Content-Type", "application/json")
}

// createRequest sends a request to the specified resource path.
func createRequest(resourcePath string) ([]byte, error) {
	req, err := http.NewRequest("GET", apiEndpointBase+resourcePath, nil)
	if err != nil {
		return []byte{}, err
	}

	// fmt.Println("fetching:", apiEndpointBase+req.URL.Path)

	sig := createSignature(req)
	appendHeaders(req, sig)

	hc := http.Client{}
	resp, err := hc.Do(req)

	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("bad HTTP status return code: %v\n%v", resp.Status, string(body))
	}

	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

//
// ───────────────────────────────────────────────────────── HELPER FUNCTIONS ─────
//
