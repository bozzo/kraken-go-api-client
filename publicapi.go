package krakenapi

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// List of valid public methods
var publicMethods = []string{
	"Assets",
	"AssetPairs",
	"Depth",
	"OHLC",
	"OHLCMinutes",
	"Spread",
	"Ticker",
	"Time",
	"Trades",
}

type PublicAPI interface {
	Time() (*TimeResponse, error)
	Assets() (AssetsResponse, error)
	AssetPairs() (AssetPairsResponse, error)
	Ticker(pairs ...string) (TickerResponse, error)
	OHLC(pair string, interval string, since int64) (*OHLCResponse, error)
	OHLCMinutes(pair string) (*OHLCResponse, error)
	Trades(pair string, since int64) (*TradesResponse, error)
	Depth(pair string, count int) (*OrderBook, error)
}

// krakenAPI represents a Kraken API Client connection
type KrakenPublic struct {
	KrakenClient
}

// Time returns the server's time
func (api *KrakenPublic) Time() (*TimeResponse, error) {
	resp, err := api.queryPublic("Time", nil, &TimeResponse{})
	if err != nil {
		return nil, err
	}

	return resp.(*TimeResponse), nil
}

// Assets returns the servers available assets
func (api *KrakenPublic) Assets() (AssetsResponse, error) {
	resp, err := api.queryPublic("Assets", nil, &Assets{})
	if err != nil {
		return nil, err
	}

	return resp.(AssetsResponse), nil
}

// AssetPairs returns the servers available asset pairs
func (api *KrakenPublic) AssetPairs() (AssetPairsResponse, error) {
	resp, err := api.queryPublic("AssetPairs", nil, &AssetPairs{})
	if err != nil {
		return nil, err
	}

	return resp.(AssetPairsResponse), nil
}

// Ticker returns the ticker for given comma separated pairs
func (api *KrakenPublic) Ticker(pairs ...string) (TickerResponse, error) {
	resp, err := api.queryPublic("Ticker", url.Values{
		"pair": {strings.Join(pairs, ",")},
	}, &Tickers{})
	if err != nil {
		return nil, err
	}

	return resp.(TickerResponse), nil
}

// OHLCWithInterval returns a OHLCResponse struct based on the given pair
func (api *KrakenPublic) OHLC(pair string, interval string, since int64) (*OHLCResponse, error) {
	urlValue := url.Values{}
	urlValue.Add("pair", pair)

	if since > 0 {
		urlValue.Add("since", fmt.Sprintf("%d", since))
	}
	if interval == "" {
		urlValue.Add("interval", "1")
	} else {
		switch interval {
		// supported values https://www.kraken.com/features/api#get-ohlc-data
		case "1", "5", "15", "30", "60", "240", "1440", "10080", "21600":
			urlValue.Add("interval", interval)
		default:
			return nil, fmt.Errorf("Unsupported value for Interval: " + interval)
		}
	}

	// Returns a map[string]interface{} as an interface{}
	interfaceResponse, err := api.queryPublic("OHLC", urlValue, nil)
	if err != nil {
		return nil, err
	}

	// Converts the interface into map[string]interface{}
	mapResponse := interfaceResponse.(map[string]interface{})
	// Extracts the list of OHLC from the map to build a slice of interfaces
	OHLCsUnstructured := mapResponse[pair].([]interface{})

	ret := new(OHLCResponse)
	for _, OHLCInterfaceSlice := range OHLCsUnstructured {
		OHLCObj, OHLCErr := NewOHLC(OHLCInterfaceSlice.([]interface{}))
		if OHLCErr != nil {
			return nil, OHLCErr
		}

		ret.OHLC = append(ret.OHLC, OHLCObj)
	}

	ret.Pair = pair
	ret.Last = int64(mapResponse["last"].(float64))

	return ret, nil
}

// OHLC returns a OHLCResponse struct based on the given pair
// Backward compatible with previous version
func (api *KrakenPublic) OHLCMinutes(pair string) (*OHLCResponse, error) {
	ret, err := api.OHLC(pair, "1", 0)

	return ret, err
}

// Trades returns the recent trades for given pair
func (api *KrakenPublic) Trades(pair string, since int64) (*TradesResponse, error) {
	values := url.Values{"pair": {pair}}
	if since > 0 {
		values.Set("since", strconv.FormatInt(since, 10))
	}
	resp, err := api.queryPublic("Trades", values, nil)
	if err != nil {
		return nil, err
	}

	v := resp.(map[string]interface{})

	last, err := strconv.ParseInt(v["last"].(string), 10, 64)
	if err != nil {
		return nil, err
	}

	result := &TradesResponse{
		Last:   last,
		Trades: make([]TradeInfo, 0),
	}

	trades := v[pair].([]interface{})
	for _, v := range trades {
		trade := v.([]interface{})

		priceString := trade[0].(string)
		price, _ := strconv.ParseFloat(priceString, 64)

		volumeString := trade[1].(string)
		volume, _ := strconv.ParseFloat(trade[1].(string), 64)

		tradeInfo := TradeInfo{
			Price:         priceString,
			PriceFloat:    price,
			Volume:        volumeString,
			VolumeFloat:   volume,
			Time:          int64(trade[2].(float64)),
			Buy:           trade[3].(string) == BUY,
			Sell:          trade[3].(string) == SELL,
			Market:        trade[4].(string) == MARKET,
			Limit:         trade[4].(string) == LIMIT,
			Miscellaneous: trade[5].(string),
		}

		result.Trades = append(result.Trades, tradeInfo)
	}

	return result, nil
}

// Depth returns the order book for given pair and orders count.
func (api *KrakenPublic) Depth(pair string, count int) (*OrderBook, error) {
	dr := DepthResponse{}
	_, err := api.queryPublic("Depth", url.Values{
		"pair": {pair}, "count": {strconv.Itoa(count)},
	}, &dr)

	if err != nil {
		return nil, err
	}

	if book, found := dr[pair]; found {
		return &book, nil
	}

	return nil, errors.New("invalid response")
}

// Execute a public method query
func (api *KrakenPublic) queryPublic(method string, values url.Values, typ interface{}) (interface{}, error) {
	apiUrl := fmt.Sprintf("%s/%s/public/%s", APIURL, APIVersion, method)
	resp, err := api.doRequest(apiUrl, values, nil, typ)

	return resp, err
}
