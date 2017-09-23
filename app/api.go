package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//{"pairing_id":1,"primary_currency":"THB","secondary_currency":"BTC","change":5.5,"last_price":117103,"volume_24hours":985.99500309,
//"orderbook":{"bids":{"total":652,"volume":1584.9082265,"highbid":116960},"asks":{"total":2001,"volume":482.70775266,"highbid":117700}}

func get_body() io.Reader {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha256.New()
	h.Write([]byte(api_key + nonce + api_secret))
	form := url.Values{
		"key":       {api_key},
		"nonce":     {nonce},
		"signature": {hex.EncodeToString(h.Sum(nil))},
	}
	return bytes.NewBufferString(form.Encode())
}

func api_usdthb() {
	Delay.Refresh_USDTHB = 3600
	//http://api.fixer.io/latest?base=USD&symbols=THB
	resp, err := http.Get("http://api.fixer.io/latest?base=USD&symbols=THB")
	if err != nil {
		_err("api_usdthb - ", err.Error())
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &USDTHB); err != nil {
		_err("api_usdthb - Unmarshal - ", err.Error(), string(body))
		return
	}
}

func api_pair() {
	Delay.Refresh_Pair = _ir(5, 2)
	resp, err := http.Get("https://bx.in.th/api/")
	if err != nil {
		_err("api_pair - ", err.Error())
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var pair map[string]APair
	if err := json.Unmarshal(body, &pair); err != nil {
		_err("api_pair - Unmarshal - ", err.Error(), string(body))
		return
	}
	for i := range pair {
		if is, err := strconv.ParseInt(i, 10, 64); err == nil {
			if _, ok := Bot[is]; ok {
				Bot[is].Pair = pair[i]
			}
		}
	}
}

func api_balance() {
	Delay.Refresh_Balance = _ir(10, 2)
	rsp, err := http.Post(api_url+"balance", "application/x-www-form-urlencoded", get_body())
	if err != nil {
		_err("api_balance - Post - ", err.Error())
		return
	}
	defer rsp.Body.Close()
	ct, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		_err("api_balance - ioutil - ", err.Error())
		return
	}
	ct = bytes.Replace(ct, []byte(`":"`), []byte(`":`), -1)
	ct = bytes.Replace(ct, []byte(`","`), []byte(`,"`), -1)
	ct = bytes.Replace(ct, []byte(`"}`), []byte(`}`), -1)
	var dat UIBalance
	if err := json.Unmarshal(ct, &dat); err != nil {
		_err("api_balance - Unmarshal - ", err.Error(), string(ct))
		return
	}
	if dat.Success == true {
		G_Balance = dat.Balance
	}
}

func api_order(pair int64) {
	Bot[pair].Delay.Refresh_Order = _ir(30, 5)
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha256.New()
	h.Write([]byte(api_key + nonce + api_secret))
	form := url.Values{
		"key":       {api_key},
		"nonce":     {nonce},
		"signature": {hex.EncodeToString(h.Sum(nil))},
		"pairing":   {_is(pair)},
	}
	rsp, err := http.Post(api_url+"getorders", "application/x-www-form-urlencoded", bytes.NewBufferString(form.Encode()))
	if err != nil {
		_err("api_order - Post - ", err.Error())
		return
	}
	defer rsp.Body.Close()
	ct, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		_err("api_order - ioutil - ", err.Error())
		return
	}
	var dat UIOrder
	if err := json.Unmarshal(ct, &dat); err != nil {
		_err("api_order - Unmarshal - ", err.Error(), string(ct))
		return
	}
	if dat.Success == true {
		Bot[pair].Order = []AOrder{}
		found := false
		min_sell := 0.0
		for i := range dat.Order {
			Bot[pair].Order = append(Bot[pair].Order, dat.Order[i])
			//if dat.Order[i].Type == "sell" {
			if !found {
				min_sell = dat.Order[i].Rate
				found = true
			} else {
				if dat.Order[i].Rate < min_sell {
					min_sell = dat.Order[i].Rate
				}
			}
		}
		if found {
			Bot[pair].Min_Sell = min_sell
		} else {
			Bot[pair].Min_Sell = 0
		}
	}
}

func api_history() {
	Delay.Refresh_History = _ir(30, 5)
	rsp, err := http.Post(api_url+"history", "application/x-www-form-urlencoded", get_body())
	if err != nil {
		_err("api_history - Post - ", err.Error())
		return
	}
	defer rsp.Body.Close()
	ct, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		_err("api_history - ioutil - ", err.Error())
		return
	}
	//transaction_id":"
	ct = bytes.Replace(ct, []byte(`transaction_id":"`), []byte(`transaction_id":`), -1)
	ct = bytes.Replace(ct, []byte(`","currency"`), []byte(`,"currency"`), -1)
	ct = bytes.Replace(ct, []byte(`"amount":"`), []byte(`"amount":`), -1)
	ct = bytes.Replace(ct, []byte(`","date"`), []byte(`,"date"`), -1)

	//fmt.Println(string(ct))

	var dat UIHistory
	if err := json.Unmarshal(ct, &dat); err != nil {
		_err("api_history - Unmarshal - ", err.Error(), string(ct))
		return
	}
	if dat.Success == true {
		for j := range Bot {
			Bot[j].Trans = []ATrans{}
		}

		line := ATrans{}
		date := ""
		idx := 0
		for i := range dat.Transaction {
			if date != dat.Transaction[i].Date {
				date = dat.Transaction[i].Date
				line = ATrans{}
				line.Date = dat.Transaction[i].Date
				idx = 1
			}
			if dat.Transaction[i].Type == "trade" {
				if line.Primary == 0 {
					line.Primary = dat.Transaction[i].Amout
					line.Primary_Currency = dat.Transaction[i].Currency
					idx++
				} else if line.Secondary == 0 {
					line.Secondary = dat.Transaction[i].Amout
					line.Secondary_Currency = dat.Transaction[i].Currency
					idx++
				}
			} else if dat.Transaction[i].Type == "fee" {
				line.Fee = dat.Transaction[i].Amout
				idx++
			}
			if idx == 4 {
				//Trans = append(Trans, line)

				for j := range Bot {
					if Bot[j].Pair.Primary == line.Primary_Currency && Bot[j].Pair.Secondary == line.Secondary_Currency {
						Bot[j].Trans = append(Bot[j].Trans, line)
						break
					} else if Bot[j].Pair.Primary == line.Secondary_Currency && Bot[j].Pair.Secondary == line.Primary_Currency {
						Bot[j].Trans = append(Bot[j].Trans, ATrans{Date: line.Date, Fee: line.Fee, Primary: line.Secondary, Secondary: line.Primary})
						break
					}
				}
			}
		}
	}
}

func api_book(pair int64) {
	Bot[pair].Delay.Refresh_Book = _ir(30, 5)
	resp, err := http.Get("https://bx.in.th/api/orderbook/?pairing=" + _is(pair))
	if err != nil {
		_err("api_book - ", err.Error())
		return
	}
	ct, _ := ioutil.ReadAll(resp.Body)
	ct = bytes.Replace(ct, []byte(`["`), []byte(`[`), -1)
	ct = bytes.Replace(ct, []byte(`","`), []byte(`,`), -1)
	ct = bytes.Replace(ct, []byte(`"]`), []byte(`]`), -1)
	//fmt.Println(string(ct))
	//  "bids":[["315.12000000","8.68529632"],["

	if err := json.Unmarshal(ct, &Bot[pair].Book); err != nil {
		_err("api_book - Unmarshal - ", err.Error())
		return
	}

	up_vol := 0.0
	down_vol := 0.0

	Bot[pair].Trend.UP_SUM_All = 0.0
	Bot[pair].Trend.DOWN_SUM_All = 0.0

	for i := range Bot[pair].Book.Bids {
		md, _ := Bot[pair].Book.Bids[i].([]interface{})
		down_vol += md[1].(float64)
		Bot[pair].Trend.DOWN_SUM_All += md[0].(float64) * md[1].(float64)
	}
	for i := range Bot[pair].Book.Asks {
		md, _ := Bot[pair].Book.Asks[i].([]interface{})
		up_vol += md[1].(float64)
		Bot[pair].Trend.UP_SUM_All += md[0].(float64) * md[1].(float64)
	}

	if down_vol > 0 && up_vol > 0 {
		Bot[pair].Trend.DOWN_AVG_All = _price(pair, Bot[pair].Trend.DOWN_SUM_All/down_vol)
		Bot[pair].Trend.UP_AVG_All = _price(pair, Bot[pair].Trend.UP_SUM_All/up_vol)
		vol := down_vol + up_vol
		Bot[pair].Trend.DOWN_Vol_All = int64((down_vol / vol) * 100)
		Bot[pair].Trend.UP_Vol_All = int64((up_vol / vol) * 100)
		Bot[pair].Trend.Price_AVG_All = _price(pair, (Bot[pair].Trend.DOWN_SUM_All+Bot[pair].Trend.UP_SUM_All)/(down_vol+up_vol))
	}
}

func api_trade(pair int64) {
	Bot[pair].Delay.Refresh_Trade = _ir(30, 5)
	resp, err := http.Get("https://bx.in.th/api/trade/?pairing=" + _is(pair) + "&limit=10")
	if err != nil {
		_err("api_trade - ", err.Error())
		return
	}
	ct, _ := ioutil.ReadAll(resp.Body)
	ct = bytes.Replace(ct, []byte(`trade_id":"`), []byte(`trade_id":`), -1)
	ct = bytes.Replace(ct, []byte(`","rate`), []byte(`,"rate`), -1)
	ct = bytes.Replace(ct, []byte(`rate":"`), []byte(`rate":`), -1)
	ct = bytes.Replace(ct, []byte(`","amount`), []byte(`,"amount`), -1)
	ct = bytes.Replace(ct, []byte(`amount":"`), []byte(`amount":`), -1)
	ct = bytes.Replace(ct, []byte(`","trade_date`), []byte(`,"trade_date`), -1)
	ct = bytes.Replace(ct, []byte(`order_id":"`), []byte(`order_id":`), -1)
	ct = bytes.Replace(ct, []byte(`","trade_type`), []byte(`,"trade_type`), -1)
	ct = bytes.Replace(ct, []byte(`","date_added`), []byte(`,"date_added`), -1)
	ct = bytes.Replace(ct, []byte(`reference_id":"`), []byte(`reference_id":`), -1)
	ct = bytes.Replace(ct, []byte(`","seconds`), []byte(`,"seconds`), -1)
	ct = bytes.Replace(ct, []byte(`seconds":"`), []byte(`seconds":`), -1)
	ct = bytes.Replace(ct, []byte(`"},{"trade_id`), []byte(`},{"trade_id`), -1)

	/*
			   "trade_id":1790148,"rate":325.00000000,"amount":1.00000000,"trade_date":2017-09-17 06:15:01,"order_id":826224,"trade_type":sell,"reference_id":0,"seconds":611},{"trade_id"

		,{"order_id":825895,"rate":326.00000000,"amount":93.18402930,"date_added":2017-09-17 03:51:33,"order_type":sell,"display_vol1":30,377.99 THB,"display_vol2":93.18402930 OMG},{"order_id"
	*/
	//fmt.Println(string(ct))
	var dat UITrade
	if err := json.Unmarshal(ct, &dat); err != nil {
		_err("api_trade - Unmarshal - ", err.Error())
		return
	}
	//fmt.Println(dat)
	trade_vol := 0.0
	up_vol := 0.0
	down_vol := 0.0

	Bot[pair].Trend.TRADE_SUM = 0.0
	Bot[pair].Trend.UP_SUM_10 = 0.0
	Bot[pair].Trend.DOWN_SUM_10 = 0.0

	/*
		ID    int32   `json:"order_id"`
		Rate  float64 `json:"rate"`
		Amout float64 `json:"amount"`
		Date  string  `json:"date_added"`
		Type  string  `json:"order_type"`
		Vol1  string  `json:"display_vol1"`
		Vol2  string  `json:"display_vol2"`
	*/
	for i := range dat.Complete {
		trade_vol += dat.Complete[i].Amout
		Bot[pair].Trend.TRADE_SUM += dat.Complete[i].Amout * dat.Complete[i].Rate
	}
	if trade_vol > 0 {
		Bot[pair].Trend.TRADE_AVG = _price(pair, Bot[pair].Trend.TRADE_SUM/trade_vol)
	}
	for i := range dat.Bids {
		down_vol += dat.Bids[i].Amout
		Bot[pair].Trend.DOWN_SUM_10 += dat.Bids[i].Amout * dat.Bids[i].Rate
	}
	for i := range dat.Asks {
		up_vol += dat.Asks[i].Amout
		Bot[pair].Trend.UP_SUM_10 += dat.Asks[i].Amout * dat.Asks[i].Rate
	}
	if down_vol > 0 && up_vol > 0 {
		Bot[pair].Trend.DOWN_AVG_10 = _price(pair, Bot[pair].Trend.DOWN_SUM_10/down_vol)
		Bot[pair].Trend.UP_AVG_10 = _price(pair, Bot[pair].Trend.UP_SUM_10/up_vol)
		vol := down_vol + up_vol
		Bot[pair].Trend.DOWN_Vol_10 = int64((down_vol / vol) * 100)
		Bot[pair].Trend.UP_Vol_10 = int64((up_vol / vol) * 100)
		Bot[pair].Trend.Price_AVG_10 = _price(pair, (Bot[pair].Trend.DOWN_SUM_10+Bot[pair].Trend.UP_SUM_10)/(down_vol+up_vol))
	}
}

func api_buy(pair int64, amount float64, rate float64) {
	if Bot[pair].Delay.Next_Buy > 0 {
		return
	}
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha256.New()
	h.Write([]byte(api_key + nonce + api_secret))
	form := url.Values{
		"key":       {api_key},
		"nonce":     {nonce},
		"signature": {hex.EncodeToString(h.Sum(nil))},
		"pairing":   {_is(pair)},
		"type":      {"buy"},
		"amount":    {_fs(amount)},
		"rate":      {_fs(rate)},
	}
	Bot[pair].Delay.Next_Buy = 60
	rsp, err := http.Post(api_url+"order", "application/x-www-form-urlencoded", bytes.NewBufferString(form.Encode()))
	if err != nil {
		_err("api_buy - Post - ", err.Error())
		return
	}
	Bot[pair].Delay.Next_Buy = 300
	Bot[pair].Delay.Next_Sell = Bot[pair].Delay.Next_Sell + 5
	defer rsp.Body.Close()
	ct, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		_err("api_buy - ioutil - ", err.Error())
		return
	}
	//	fmt.Println(string(ct))
	var dat UIOrder
	if err := json.Unmarshal(ct, &dat); err != nil {
		_err("api_buy - Unmarshal - ", err.Error())
		return
	}
	Bot[pair].Delay.Next_Sell = Bot[pair].Delay.Next_Sell + 5
	Bot[pair].Delay.Refresh_Order = 0
	if dat.Success == true {

	}
}

func api_sell(pair int64, amount float64, rate float64) {
	if Bot[pair].Delay.Next_Sell > 0 {
		return
	}
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha256.New()
	h.Write([]byte(api_key + nonce + api_secret))
	form := url.Values{
		"key":       {api_key},
		"nonce":     {nonce},
		"signature": {hex.EncodeToString(h.Sum(nil))},
		"pairing":   {_is(pair)},
		"type":      {"sell"},
		"amount":    {_fs(amount)},
		"rate":      {_fs(rate)},
	}
	Bot[pair].Delay.Next_Sell = 60
	Bot[pair].Delay.Next_Buy = 60
	rsp, err := http.Post(api_url+"order", "application/x-www-form-urlencoded", bytes.NewBufferString(form.Encode()))
	if err != nil {
		_err("api_sell - Post - ", err.Error())
		return
	}
	Bot[pair].Delay.Next_Sell = 300
	defer rsp.Body.Close()
	ct, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		_err("api_sell - ioutil - ", err.Error())
		return
	}
	//	fmt.Println(string(ct))
	var dat UIOrder
	if err := json.Unmarshal(ct, &dat); err != nil {
		_err("api_sell - Unmarshal - ", err.Error())
		return
	}
	Bot[pair].Delay.Next_Buy = 300
	Bot[pair].Delay.Refresh_Order = 0
	if dat.Success == true {

	}
}

func api_cancel(pair int64, id int64) {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha256.New()
	h.Write([]byte(api_key + nonce + api_secret))
	form := url.Values{
		"key":       {api_key},
		"nonce":     {nonce},
		"signature": {hex.EncodeToString(h.Sum(nil))},
		"pairing":   {_is(pair)},
		"order_id":  {_is(id)},
	}
	rsp, err := http.Post(api_url+"cancel", "application/x-www-form-urlencoded", bytes.NewBufferString(form.Encode()))
	if err != nil {
		_err("api_cancel - Post - ", err.Error())
		return
	}
	defer rsp.Body.Close()
	ct, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		_err("api_cancel - ioutil - ", err.Error())
		return
	}
	//	fmt.Println(string(ct))
	var dat UIOrder
	if err := json.Unmarshal(ct, &dat); err != nil {
		_err("api_cancel - Unmarshal - ", err.Error())
		return
	}
	for i := range Bot[pair].Order {
		if Bot[pair].Order[i].ID == id {
			Bot[pair].Order = append(Bot[pair].Order[:i], Bot[pair].Order[i+1:]...)
			break
		}
	}
	Bot[pair].Delay.Refresh_Order = 0
	if dat.Success == true {

	}
}
