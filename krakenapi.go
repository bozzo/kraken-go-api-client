package krakenapi

import (
	"net/http"
)

const (
	// APIURL is the official Kraken API Endpoint
	APIURL = "https://api.kraken.com"
	// APIVersion is the official Kraken API Version Number
	APIVersion = "0"
)

// KrakenApi represents a Kraken API Client connection
type API interface {
	Public() PublicAPI
	Private() PrivateAPI
}

// krakenAPI represents a Kraken API Client connection
type krakenAPI struct {
	public  PublicAPI
	private PrivateAPI
}

// New creates a new Kraken API client
func New(key, secret string) API {
	return NewWithClient(key, secret, http.DefaultClient)
}

// NewWithClient creates a new Kraken API client with custom http client
func NewWithClient(key, secret string, httpClient *http.Client) API {
	api := &krakenAPI{
		public: &KrakenPublic{
			KrakenClient{
				client: httpClient,
			},
		},
		private: &KrakenPrivate{
			key:    key,
			secret: secret,
			KrakenClient: KrakenClient{
				client: httpClient,
			},
		},
	}
	return api
}

func (api *krakenAPI) Public() PublicAPI {
	return api.public
}

func (api *krakenAPI) Private() PrivateAPI {
	return api.private
}