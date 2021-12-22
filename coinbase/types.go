package coinbase

import (
	"time"
)

var (
	cbAPIKey        string
	cbAPISecret     string
	cbAPIVersion    string = "2017-08-31"
	apiEndpointBase string = "https://api.coinbase.com/v2/"
)

// These constants are used to map the types of prices that can be used to pass to the
// GetPrice() method.
const (
	Buy             string = "buy"
	Sell            string = "sell"
	Spot            string = "spot"
	InflationReward string = "inflation_reward"
)

type CoinbaseClient struct{}

// User is a structure containing user profile information parsed from the https://api.coinbase.com/v2/user api endpoint path.
type User struct {
	Data struct {
		ID              string      `json:"id"`
		Name            string      `json:"name"`
		Username        interface{} `json:"username"`
		ProfileLocation interface{} `json:"profile_location"`
		ProfileBio      interface{} `json:"profile_bio"`
		ProfileURL      interface{} `json:"profile_url"`
		AvatarURL       string      `json:"avatar_url"`
		Resource        string      `json:"resource"`
		ResourcePath    string      `json:"resource_path"`
		LegacyID        string      `json:"legacy_id"`
		TimeZone        string      `json:"time_zone"`
		NativeCurrency  string      `json:"native_currency"`
		BitcoinUnit     string      `json:"bitcoin_unit"`
		State           string      `json:"state"`
		Country         struct {
			Code       string `json:"code"`
			Name       string `json:"name"`
			IsInEurope bool   `json:"is_in_europe"`
		} `json:"country"`
		Nationality struct {
			Code interface{} `json:"code"`
			Name interface{} `json:"name"`
		} `json:"nationality"`
		RegionSupportsFiatTransfers           bool      `json:"region_supports_fiat_transfers"`
		RegionSupportsCryptoToCryptoTransfers bool      `json:"region_supports_crypto_to_crypto_transfers"`
		CreatedAt                             time.Time `json:"created_at"`
		SupportsRewards                       bool      `json:"supports_rewards"`
		Tiers                                 struct {
			CompletedDescription string      `json:"completed_description"`
			UpgradeButtonText    interface{} `json:"upgrade_button_text"`
			Header               interface{} `json:"header"`
			Body                 interface{} `json:"body"`
		} `json:"tiers"`
		ReferralMoney struct {
			Amount            string `json:"amount"`
			Currency          string `json:"currency"`
			CurrencySymbol    string `json:"currency_symbol"`
			ReferralThreshold string `json:"referral_threshold"`
		} `json:"referral_money"`
		HasBlockingBuyRestrictions            bool   `json:"has_blocking_buy_restrictions"`
		HasMadeAPurchase                      bool   `json:"has_made_a_purchase"`
		HasBuyDepositPaymentMethods           bool   `json:"has_buy_deposit_payment_methods"`
		HasUnverifiedBuyDepositPaymentMethods bool   `json:"has_unverified_buy_deposit_payment_methods"`
		NeedsKycRemediation                   bool   `json:"needs_kyc_remediation"`
		ShowInstantAchUx                      bool   `json:"show_instant_ach_ux"`
		UserType                              string `json:"user_type"`
	} `json:"data"`
}

// Account is a structure containing account information parsed from the https://api.coinbase.com/v2/accounts api endpoint path.
type Account struct {
	Pagination struct {
		EndingBefore  interface{} `json:"ending_before"`
		StartingAfter interface{} `json:"starting_after"`
		Limit         int         `json:"limit"`
		Order         string      `json:"order"`
		PreviousURI   interface{} `json:"previous_uri"`
		NextURI       interface{} `json:"next_uri"`
	} `json:"pagination"`
	Data []struct {
		ID       string      `json:"id"`
		Name     string      `json:"name"`
		Primary  bool        `json:"primary"`
		Type     string      `json:"type"`
		Currency interface{} `json:"currency"`
		Balance  struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"balance"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Resource     string    `json:"resource"`
		ResourcePath string    `json:"resource_path"`
		Ready        bool      `json:"ready,omitempty"`
	} `json:"data"`
}

// ExchangeRate is used to parse the current exchange rates for crypto currencies available in Coinbase.
type ExchangeRate map[string]interface{}

// Price is used to parse the current spot price for a specified crypto currency.
type Price struct {
	Data struct {
		Base     string `json:"base"`
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"data"`
}

// Transaction is used to parse the transaction history of a specified account.
type Transaction struct {
	Data []struct {
		ID     string `json:"id"`
		Type   string `json:"type"`
		Status string `json:"status"`
		Amount struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"amount"`
		NativeAmount struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"native_amount"`
		Description     interface{} `json:"description"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		Resource        string      `json:"resource"`
		ResourcePath    string      `json:"resource_path"`
		InstantExchange bool        `json:"instant_exchange"`
		Buy             struct {
			ID           string `json:"id"`
			Resource     string `json:"resource"`
			ResourcePath string `json:"resource_path"`
		} `json:"buy"`
		Details struct {
			Title             string `json:"title"`
			Subtitle          string `json:"subtitle"`
			Header            string `json:"header"`
			Health            string `json:"health"`
			PaymentMethodName string `json:"payment_method_name"`
		} `json:"details"`
		HideNativeAmount bool `json:"hide_native_amount"`
	} `json:"data"`
	Pagination struct {
		EndingBefore         interface{} `json:"ending_before"`
		StartingAfter        interface{} `json:"starting_after"`
		PreviousEndingBefore interface{} `json:"previous_ending_before"`
		NextStartingAfter    interface{} `json:"next_starting_after"`
		Limit                int         `json:"limit"`
		Order                string      `json:"order"`
		PreviousURI          interface{} `json:"previous_uri"`
		NextURI              interface{} `json:"next_uri"`
	} `json:"pagination"`
}
