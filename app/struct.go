package main

var api_url = "https://bx.in.th/api/"
var api_key = ""
var api_secret = ""
var api_pass = ""
var api_port int64 = 0

// global
var G_Balance map[string]ABalance
var Bot map[int64]*ABot
var Delay GDelay
var VERSION string
var Sort []int64
var Bitfinex map[string][]float64
var USDTHB AUSDTHB
var Conf GConf

// each
type ABot struct {
	Pair  APair
	Conf  AConfig
	Order []AOrder
	//	History  []AHistory
	Book     UIBook
	Trade    UITrade
	Trend    ATrend
	Trans    []ATrans
	Sims     []ASims
	Delay    ADelay
	Min_Sell float64
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
	Theme string
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

/*
{"success":true,"orders":[{"pairing_id":26,"order_id":813197,"order_type":"sell","rate":347,"amount":10,"date":"2017-09-16 13:54:27"}]}
*/

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
}

type ADelay struct {
	Next_Buy      int64
	Next_Sell     int64
	Refresh_Book  int64
	Refresh_Trade int64
	Refresh_Order int64
}

/*
{"bids":[["315.12000000","8.68529632"],["315.11100000","185.71348644"],["315.11000000","21.10140959"],["315.01000000","11.74435050"],["315.00110000","15.83327804"],["315.00000000","633.33333333"],[
*/

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
