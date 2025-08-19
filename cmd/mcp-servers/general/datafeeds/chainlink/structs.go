package chainlink

type Root struct {
	Props Props `json:"props"`
}

type Props struct {
	PageProps PageProps `json:"pageProps"`
}

type PageProps struct {
	AllFeeds []Feed `json:"allFeeds"`
}

type Feed struct {
	CompareOffchain     string      `json:"compareOffchain"`
	ContractAddress     string      `json:"contractAddress"`
	ContractType        string      `json:"contractType"`
	ContractVersion     int         `json:"contractVersion"`
	DecimalPlaces       *int        `json:"decimalPlaces"`
	Ens                 string      `json:"ens"`
	FormatDecimalPlaces *int        `json:"formatDecimalPlaces"`
	HealthPrice         string      `json:"healthPrice"`
	Heartbeat           int         `json:"heartbeat"`
	History             interface{} `json:"history"`
	Multiply            string      `json:"multiply"`
	Name                string      `json:"name"`
	Pair                []string    `json:"pair"`
	Path                string      `json:"path"`
	ProxyAddress        string      `json:"proxyAddress"`
	Threshold           float64     `json:"threshold"`
	ValuePrefix         string      `json:"valuePrefix"`
	AssetName           string      `json:"assetName"`
	FeedCategory        string      `json:"feedCategory"`
	FeedType            string      `json:"feedType"`
	Docs                Docs        `json:"docs"`
	Decimals            int         `json:"decimals"`
	Symbol              Symbol      `json:"symbol"`
	Chain               string      `json:"chain"`
	Network             string      `json:"network"`
	ChainID             string      `json:"chainId"`
	AssetClass          string      `json:"assetClass"`
	ProductType         string      `json:"productType"`
	ProductTypeCode     string      `json:"productTypeCode"`
	ProductSubType      string      `json:"productSubType"`
	SvrEnabled          bool        `json:"svrEnabled"`
	SvrFeedConfig       *FeedConfig `json:"svrFeedConfig"`
	Icon                string      `json:"icon"`
	PublishedAt         interface{} `json:"publishedAt"`
	MarketData          MarketData  `json:"marketData"`
}

type Docs struct {
	AssetClass          string `json:"assetClass"`
	BaseAsset           string `json:"baseAsset"`
	BlockchainName      string `json:"blockchainName"`
	ClicProductName     string `json:"clicProductName"`
	DeliveryChannelCode string `json:"deliveryChannelCode"`
	MarketHours         string `json:"marketHours"`
	ProductSubType      string `json:"productSubType"`
	ProductType         string `json:"productType"`
	ProductTypeCode     string `json:"productTypeCode"`
	QuoteAsset          string `json:"quoteAsset"`
	QuoteAssetClic      string `json:"quoteAssetClic"`
	Hidden              bool   `json:"hidden,omitempty"`
}

type Symbol struct {
	Prefix string `json:"prefix"`
}

type FeedConfig struct {
	CompareOffchain       string      `json:"compareOffchain"`
	ContractAddress       string      `json:"contractAddress"`
	ContractType          string      `json:"contractType"`
	ContractVersion       int         `json:"contractVersion"`
	DecimalPlaces         *int        `json:"decimalPlaces"`
	Ens                   *string     `json:"ens"`
	FormatDecimalPlaces   *int        `json:"formatDecimalPlaces"`
	HealthPrice           string      `json:"healthPrice"`
	Heartbeat             int         `json:"heartbeat"`
	History               interface{} `json:"history"`
	Multiply              string      `json:"multiply"`
	Name                  string      `json:"name"`
	Pair                  []string    `json:"pair"`
	Path                  string      `json:"path"`
	ProxyAddress          string      `json:"proxyAddress"`
	SecondaryProxyAddress string      `json:"secondaryProxyAddress"`
	Threshold             float64     `json:"threshold"`
	ValuePrefix           string      `json:"valuePrefix"`
	AssetName             string      `json:"assetName"`
	FeedCategory          string      `json:"feedCategory"`
	FeedType              string      `json:"feedType"`
	Docs                  Docs        `json:"docs"`
	Decimals              int         `json:"decimals"`
}

type MarketData struct {
	Asset             string `json:"asset"`
	MarketCap         string `json:"marketCap"`
	CirculatingSupply string `json:"circulatingSupply"`
	TotalVolume       string `json:"totalVolume"`
	Quote             string `json:"quote"`
}
