package krakenapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"time"
)

// List of valid private methods
var privateMethods = []string{
	"AddExport",
	"AddOrder",
	"Balance",
	"BalanceMap",
	"CancelOrder",
	"ClosedOrders",
	"DepositAddresses",
	"DepositMethods",
	"DepositStatus",
	"ExportStatus",
	"GetWebSocketsToken",
	"Ledgers",
	"OpenOrders",
	"OpenPositions",
	"QueryLedgers",
	"QueryOrders",
	"QueryTrades",
	"RemoveExport",
	"RetrieveExport",
	"TradeBalance",
	"TradesHistory",
	"TradeVolume",
	"WalletTransfer",
	"Withdraw",
	"WithdrawCancel",
	"WithdrawInfo",
	"WithdrawStatus",
}

type PrivateAPI interface {
	TradesHistory(start int64, end int64, args map[string]string) (*TradesHistoryResponse, error)
	Balance() (BalanceResponse, error)
	TradeBalance(args map[string]string) (*TradeBalanceResponse, error)
	TradeVolume(args map[string]string) (*TradeVolumeResponse, error)
	OpenOrders(args map[string]string) (*OpenOrdersResponse, error)
	ClosedOrders(args map[string]string) (*ClosedOrdersResponse, error)
	CancelOrder(txid string) (*CancelOrderResponse, error)
	QueryOrders(txids string, args map[string]string) (*QueryOrdersResponse, error)
	AddOrder(pair string, direction string, orderType string, volume string, args map[string]string) (*AddOrderResponse, error)
	Ledgers(args map[string]string) (*LedgersResponse, error)
	DepositAddresses(asset string, method string) (*DepositAddressesResponse, error)
	Withdraw(asset string, key string, amount *big.Float) (*WithdrawResponse, error)
	WithdrawInfo(asset string, key string, amount *big.Float) (*WithdrawInfoResponse, error)
}

// krakenAPI represents a Kraken API Client connection
type KrakenPrivate struct {
	key    string
	secret string
	KrakenClient
}

// TradesHistory returns the Trades History within a specified time frame (start to end).
func (api *KrakenPrivate) TradesHistory(start int64, end int64, args map[string]string) (*TradesHistoryResponse, error) {
	params := url.Values{}
	if start > 0 {
		params.Add("start", strconv.FormatInt(start, 10))
	}
	if end > 0 {
		params.Add("end", strconv.FormatInt(end, 10))
	}
	if value, ok := args["type"]; ok {
		params.Add("type", value)
	}
	if value, ok := args["trades"]; ok {
		params.Add("trades", value)
	}
	if value, ok := args["ofs"]; ok {
		params.Add("ofs", value)
	}

	resp, err := api.queryPrivate("TradesHistory", params, &TradesHistoryResponse{})

	if err != nil {
		return nil, err
	}

	return resp.(*TradesHistoryResponse), nil
}

// Balance returns all account asset balances
func (api *KrakenPrivate) Balance() (BalanceResponse, error) {
	resp, err := api.queryPrivate("Balance", url.Values{}, &map[string]string{})
	if err != nil {
		return nil, err
	}

	balances := Balances{}
	for key, value := range *resp.(*map[string]string) {
		balances[key], err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
	}
	return &balances, nil
}

// TradeBalance returns trade balance info
func (api *KrakenPrivate) TradeBalance(args map[string]string) (*TradeBalanceResponse, error) {
	params := url.Values{}
	if value, ok := args["aclass"]; ok {
		params.Add("aclass", value)
	}
	if value, ok := args["asset"]; ok {
		params.Add("asset", value)
	}
	resp, err := api.queryPrivate("TradeBalance", params, &TradeBalanceResponse{})
	if err != nil {
		return nil, err
	}

	return resp.(*TradeBalanceResponse), nil
}

// TradeVolume returns trade volume info
func (api *KrakenPrivate) TradeVolume(args map[string]string) (*TradeVolumeResponse, error) {
	params := url.Values{}
	if value, ok := args["pair"]; ok {
		params.Add("pair", value)
	}
	if value, ok := args["fee-info"]; ok {
		params.Add("fee-info", value)
	}
	resp, err := api.queryPrivate("TradeVolume", params, &TradeVolumeResponse{})
	if err != nil {
		return nil, err
	}

	return resp.(*TradeVolumeResponse), nil
}

// OpenOrders returns all open orders
func (api *KrakenPrivate) OpenOrders(args map[string]string) (*OpenOrdersResponse, error) {
	params := url.Values{}
	if value, ok := args["trades"]; ok {
		params.Add("trades", value)
	}
	if value, ok := args["userref"]; ok {
		params.Add("userref", value)
	}

	resp, err := api.queryPrivate("OpenOrders", params, &OpenOrdersResponse{})

	if err != nil {
		return nil, err
	}

	return resp.(*OpenOrdersResponse), nil
}

// ClosedOrders returns all closed orders
func (api *KrakenPrivate) ClosedOrders(args map[string]string) (*ClosedOrdersResponse, error) {
	params := url.Values{}
	if value, ok := args["trades"]; ok {
		params.Add("trades", value)
	}
	if value, ok := args["userref"]; ok {
		params.Add("userref", value)
	}
	if value, ok := args["start"]; ok {
		params.Add("start", value)
	}
	if value, ok := args["end"]; ok {
		params.Add("end", value)
	}
	if value, ok := args["ofs"]; ok {
		params.Add("ofs", value)
	}
	if value, ok := args["closetime"]; ok {
		params.Add("closetime", value)
	}
	resp, err := api.queryPrivate("ClosedOrders", params, &ClosedOrdersResponse{})

	if err != nil {
		return nil, err
	}

	return resp.(*ClosedOrdersResponse), nil
}

// CancelOrder cancels order
func (api *KrakenPrivate) CancelOrder(txid string) (*CancelOrderResponse, error) {
	params := url.Values{}
	params.Add("txid", txid)
	resp, err := api.queryPrivate("CancelOrder", params, &CancelOrderResponse{})

	if err != nil {
		return nil, err
	}

	return resp.(*CancelOrderResponse), nil
}

// QueryOrders shows order
func (api *KrakenPrivate) QueryOrders(txids string, args map[string]string) (*QueryOrdersResponse, error) {
	params := url.Values{"txid": {txids}}
	if value, ok := args["trades"]; ok {
		params.Add("trades", value)
	}
	if value, ok := args["userref"]; ok {
		params.Add("userref", value)
	}
	resp, err := api.queryPrivate("QueryOrders", params, &QueryOrdersResponse{})

	if err != nil {
		return nil, err
	}

	return resp.(*QueryOrdersResponse), nil
}

// AddOrder adds new order
func (api *KrakenPrivate) AddOrder(pair string, direction string, orderType string, volume string, args map[string]string) (*AddOrderResponse, error) {
	params := url.Values{
		"pair":      {pair},
		"type":      {direction},
		"ordertype": {orderType},
		"volume":    {volume},
	}

	if value, ok := args["price"]; ok {
		params.Add("price", value)
	}
	if value, ok := args["price2"]; ok {
		params.Add("price2", value)
	}
	if value, ok := args["leverage"]; ok {
		params.Add("leverage", value)
	}
	if value, ok := args["oflags"]; ok {
		params.Add("oflags", value)
	}
	if value, ok := args["starttm"]; ok {
		params.Add("starttm", value)
	}
	if value, ok := args["expiretm"]; ok {
		params.Add("expiretm", value)
	}
	if value, ok := args["validate"]; ok {
		params.Add("validate", value)
	}
	if value, ok := args["close_order_type"]; ok {
		params.Add("close[ordertype]", value)
	}
	if value, ok := args["close_price"]; ok {
		params.Add("close[price]", value)
	}
	if value, ok := args["close_price2"]; ok {
		params.Add("close[price2]", value)
	}
	if value, ok := args["trading_agreement"]; ok {
		params.Add("trading_agreement", value)
	}
	if value, ok := args["userref"]; ok {
		params.Add("userref", value)
	}
	resp, err := api.queryPrivate("AddOrder", params, &AddOrderResponse{})

	if err != nil {
		return nil, err
	}

	return resp.(*AddOrderResponse), nil
}

// Ledgers returns ledgers informations
func (api *KrakenPrivate) Ledgers(args map[string]string) (*LedgersResponse, error) {
	params := url.Values{}
	if value, ok := args["aclass"]; ok {
		params.Add("aclass", value)
	}
	if value, ok := args["asset"]; ok {
		params.Add("asset", value)
	}
	if value, ok := args["type"]; ok {
		params.Add("type", value)
	}
	if value, ok := args["start"]; ok {
		params.Add("start", value)
	}
	if value, ok := args["end"]; ok {
		params.Add("end", value)
	}
	if value, ok := args["ofs"]; ok {
		params.Add("ofs", value)
	}
	resp, err := api.queryPrivate("Ledgers", params, &LedgersResponse{})
	if err != nil {
		return nil, err
	}

	return resp.(*LedgersResponse), nil
}

// DepositAddresses returns deposit addresses
func (api *KrakenPrivate) DepositAddresses(asset string, method string) (*DepositAddressesResponse, error) {
	resp, err := api.queryPrivate("DepositAddresses", url.Values{
		"asset":  {asset},
		"method": {method},
	}, &DepositAddressesResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*DepositAddressesResponse), nil
}

// Withdraw executes a withdrawal, returning a reference ID
func (api *KrakenPrivate) Withdraw(asset string, key string, amount *big.Float) (*WithdrawResponse, error) {
	resp, err := api.queryPrivate("Withdraw", url.Values{
		"asset":  {asset},
		"key":    {key},
		"amount": {amount.String()},
	}, &WithdrawResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*WithdrawResponse), nil
}

// WithdrawInfo returns withdrawal information
func (api *KrakenPrivate) WithdrawInfo(asset string, key string, amount *big.Float) (*WithdrawInfoResponse, error) {
	resp, err := api.queryPrivate("WithdrawInfo", url.Values{
		"asset":  {asset},
		"key":    {key},
		"amount": {amount.String()},
	}, &WithdrawInfoResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*WithdrawInfoResponse), nil
}

// queryPrivate executes a private method query
func (api *KrakenPrivate) queryPrivate(method string, values url.Values, typ interface{}) (interface{}, error) {
	urlPath := fmt.Sprintf("/%s/private/%s", APIVersion, method)
	reqURL := fmt.Sprintf("%s%s", APIURL, urlPath)
	secret, _ := base64.StdEncoding.DecodeString(api.secret)
	values.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

	// Create signature
	signature := createSignature(urlPath, values, secret)

	// Add Key and signature to request headers
	headers := map[string]string{
		"API-Key":  api.key,
		"API-Sign": signature,
	}

	resp, err := api.doRequest(reqURL, values, headers, typ)

	return resp, err
}

// getSha256 creates a sha256 hash for given []byte
func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}

// getHMacSha512 creates a hmac hash with sha512
func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func createSignature(urlPath string, values url.Values, secret []byte) string {
	// See https://www.kraken.com/help/api#general-usage for more information
	shaSum := getSha256([]byte(values.Get("nonce") + values.Encode()))
	macSum := getHMacSha512(append([]byte(urlPath), shaSum...), secret)
	return base64.StdEncoding.EncodeToString(macSum)
}
