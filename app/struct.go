package main

// global
var Balance map[string]ABalance
var Bot map[int64]*ABot
var Delay GDelay
var VERSION string
var Sort []int64
var Bitfinex map[string][]float64
var Bittrex map[string]GBittrex
var USDTHB AUSDTHB
var Conf GConf
var Tmp map[int64]*ABot

// each
type ABot struct {
	Pair     APair
	Conf     AConfig
	Order    []AOrder
	Book     UIBook
	Trade    UITrade
	Trend    ATrend
	Trans    []ATrans
	Sims     []ASims
	Delay    ADelay
	Min_Sell float64
	Graph    AGraph
}

type AUSDTHB struct {
	ID   string `json:"base"`
	Date string `json:"date"`
	Rate struct {
		THB float64 `json:"THB"`
	} `json:"rates"`
}

type ACoin struct {
	BCH float64
	BTC float64
	DAS float64
	DOG float64
	ETH float64
	FTC float64
	GNO float64
	HYP float64
	LTC float64
	NMC float64
	OMG float64
	PND float64
	PPC float64
	QRK float64
	REP float64
	THB float64
	XCN float64
	XPM float64
	XPY float64
	XRP float64
	ZEC float64
}

type AConfig struct {
	Enable    bool
	Budget    float64
	Max_Price float64
	Max_Order int64
	Cycle     float64
	Margin    float64
	KEEP      map[string]float64
}

type GConf struct {
	Theme  string
	Line   string
	URL    string
	Key    string
	Secret string
	Pass   string
	TwoFA  string
	Port   int64
}

type ATrend struct {
	UP_Price_All   float64
	UP_Vol_All     int64
	UP_AVG_All     float64
	UP_SUM_All     float64
	DOWN_Price_All float64
	DOWN_Vol_All   int64
	DOWN_AVG_All   float64
	DOWN_SUM_All   float64
	Price_AVG_All  float64

	TRADE_SUM float64
	TRADE_AVG float64

	UP_Price_10   float64
	UP_Vol_10     int64
	UP_AVG_10     float64
	UP_SUM_10     float64
	DOWN_Price_10 float64
	DOWN_Vol_10   int64
	DOWN_AVG_10   float64
	DOWN_SUM_10   float64
	Price_AVG_10  float64
}

type AGraph struct {
	BX            []float64
	BX_Last       float64
	BX_Time       string
	Bittrex       []float64
	Bittrex_Last  float64
	Bittrex_Time  string
	Bitfinex      []float64
	Bitfinex_Last float64
	Bitfinex_Time string
}

type AList struct {
	ID        int64   `json:"id"`
	Primary   string  `json:"primary"`
	Secondary string  `json:"secondary"`
	Change    float64 `json:"change"`
	Price     float64 `json:"price"`
}

type APair struct {
	ID        int64   `json:"pairing_id"`
	Primary   string  `json:"primary_currency"`
	Secondary string  `json:"secondary_currency"`
	Change    float64 `json:"change"`
	Price     float64 `json:"last_price"`
	Volume    float64 `json:"volume_24hours"`
	Order     struct {
		Bids struct {
			Total  int64   `json:"total"`
			Volume float64 `json:"volume"`
			Price  float64 `json:"highbid"`
		} `json:"bids"`
		Asks struct {
			Total  int64   `json:"total"`
			Volume float64 `json:"volume"`
			Price  float64 `json:"highbid"`
		} `json:"asks"`
	} `json:"orderbook"`
}

type ABalance struct {
	Total       float64 `json:"total,int32,omitempty"`
	Available   float64 `json:"available,int32,omitempty"`
	Orders      float64 `json:"orders,int32,omitempty"`
	Withdrawals float64 `json:"withdrawals,int32,omitempty"`
	Deposits    float64 `json:"deposits,int32,omitempty"`
	Options     float64 `json:"options,int32,omitempty"`
}

type UIBalance struct {
	Success bool                `json:"success"`
	Balance map[string]ABalance `json:"balance"`
	Error   string              `json:"error"`
}

type AOrder struct {
	ID     int64   `json:"order_id"`
	Pair   int64   `json:"pairing_id"`
	Type   string  `json:"order_type"`
	Rate   float64 `json:"rate"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
}

type UIOrder struct {
	Success bool     `json:"success"`
	Order   []AOrder `json:"orders"`
	Error   string   `json:"error"`
}

type UIBook struct {
	Bids []interface{} `json:"bids"`
	Asks []interface{} `json:"asks"`
}

type Atrade_Complete struct {
	ID        int64   `json:"trade_id"`
	Rate      float64 `json:"rate"`
	Amout     float64 `json:"amount"`
	Date      string  `json:"trade_date"`
	Order     int64   `json:"order_id"`
	Type      string  `json:"trade_type"`
	Reference int64   `json:"reference_id"`
	Seconds   int64   `json:"seconds"`
}

type Atrade_Order struct {
	ID    int64   `json:"order_id"`
	Rate  float64 `json:"rate"`
	Amout float64 `json:"amount"`
	Date  string  `json:"date_added"`
	Type  string  `json:"order_type"`
	Vol1  string  `json:"display_vol1"`
	Vol2  string  `json:"display_vol2"`
}

type UITrade struct {
	Complete []Atrade_Complete `json:"trades"`
	Asks     []Atrade_Order    `json:"lowask"`
	Bids     []Atrade_Order    `json:"highbid"`
}

type AHistory struct {
	ID       int64   `json:"transaction_id"`
	Currency string  `json:"currency"`
	Amout    float64 `json:"amount"`
	Date     string  `json:"date"`
	Type     string  `json:"type"`
}
type UIHistory struct {
	Success     bool       `json:"success"`
	Transaction []AHistory `json:"transactions"`
	Error       string     `json:"error"`
}

type ATrans struct {
	Date               string
	Primary            float64
	Primary_Currency   string
	Secondary          float64
	Secondary_Currency string
	Fee                float64
}

type GDelay struct {
	Refresh_Pair    int64
	Refresh_Balance int64
	Refresh_History int64
	Refresh_USDTHB  int64
	Refresh_Bittrex int64
	Next_BuySell    int64
}

type ADelay struct {
	Next_Buy      int64
	Next_Sell     int64
	Refresh_Book  int64
	Refresh_Trade int64
	Refresh_Order int64
}

type ASims struct {
	Buy        float64
	Sell       float64
	Order_Sell float64
	Order_Buy  float64
	Margin     float64
	Coin       float64
	Profit     float64
	Diff       float64
}

type GBittrex struct {
	Price      float64 `json:"Price"`
	Change     float64 `json:"Change"`
	Volume     float64 `json:"Volume"`
	Bid        float64 `json:"Bid"`
	Ask        float64 `json:"Ask"`
	Order_Buy  int64   `json:"Order_Buy"`
	Order_Sell int64   `json:"Order_Sell"`
}

type ABittrex struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  []struct {
		Market     string  `json:"MarketName"`
		High       float64 `json:"High"`
		Low        float64 `json:"Low"`
		Volume     float64 `json:"Volume"`
		Price      float64 `json:"Last"`
		BaseVolume float64 `json:"BaseVolume"`
		Date       string  `json:"TimeStamp"`
		Bid        float64 `json:"Bid"`
		Ask        float64 `json:"Ask"`
		Order_Buy  int64   `json:"OpenBuyOrders"`
		Order_Sell int64   `json:"OpenSellOrders"`
		PrevDay    float64 `json:"PrevDay"`
		Created    string  `json:"Created"`
	} `json:"result"`
}
