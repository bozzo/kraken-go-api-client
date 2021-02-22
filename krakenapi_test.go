package krakenapi

import (
	"reflect"
	"testing"
)

var api = New("", "")

func TestTime(t *testing.T) {
	resp, err := api.Public().Time()
	if err != nil {
		t.Errorf("Time() should not return an error, got %s", err)
	}

	if resp.Unixtime <= 0 {
		t.Errorf("Time() should return valid Unixtime, got %d", resp.Unixtime)
	}
}

func TestAssets(t *testing.T) {
	_, err := api.Public().Assets()
	if err != nil {
		t.Errorf("Assets() should not return an error, got %s", err)
	}
}

func TestAssetPairs(t *testing.T) {
	resp, err := api.Public().AssetPairs()
	if err != nil {
		t.Errorf("AssetPairs() should not return an error, got %s", err)
	}

	if resp.GetAssetPair(XXBTZEUR).Base+resp.GetAssetPair(XXBTZEUR).Quote != XXBTZEUR {
		t.Errorf("AssetPairs() should return valid response, got %+v", resp.GetAssetPair(XXBTZEUR))
	}

	if len(resp.GetPairs()) <= 0 {
		t.Errorf("AssetPairs GetPairs() should return the pair list, got %+v", resp.GetPairs())
	}
}

func TestTicker(t *testing.T) {
	resp, err := api.Public().Ticker(XXBTZEUR, XXRPZEUR)
	if err != nil {
		t.Errorf("Ticker() should not return an error, got %s", err)
	}

	if resp.GetPairTickerInfo(XXBTZEUR).OpeningPrice == 0 {
		t.Errorf("Ticker() should return valid OpeningPrice, got %+v", resp.GetPairTickerInfo(XXBTZEUR).OpeningPrice)
	}
}

func TestOHLC(t *testing.T) {
	resp, err := api.Public().OHLC(XXBTZEUR)
	if err != nil {
		t.Errorf("OHLC() should not return an error, got %s", err)
	}

	if resp.Pair == "" {
		t.Errorf("OHLC() should return valid Pair, got %+v", resp.Pair)
	}
}

func TestQueryTrades(t *testing.T) {
	result, err := api.Public().Trades(XXBTZEUR, 1495777604391411290)

	if err != nil {
		t.Errorf("Trades should not return an error, got %s", err)
	}

	if result.Last == 0 {
		t.Errorf("Returned parameter `last` should always have a value...")
	}

	if len(result.Trades) > 0 {
		for _, trade := range result.Trades {
			if trade.Buy == trade.Sell {
				t.Errorf("Trade should be buy or sell")
			}
			if trade.Market == trade.Limit {
				t.Errorf("Trade type should be market or limit")
			}
		}
	}
}

func TestQueryDepth(t *testing.T) {
	pair := "XETHZEUR"
	count := 10
	result, err := api.Public().Depth(pair, count)
	if err != nil {
		t.Errorf("Depth should not return an error, got %s", err)
	}

	resultType := reflect.TypeOf(result)

	if resultType != reflect.TypeOf(&OrderBook{}) {
		t.Errorf("Depth should return an OrderBook, got %s", resultType)
	}

	if len(result.Asks) > count {
		t.Errorf("Asks length must be less than count , got %d > %d", len(result.Asks), count)
	}

	if len(result.Bids) > count {
		t.Errorf("Bids length must be less than count , got %d > %d", len(result.Bids), count)
	}
}
