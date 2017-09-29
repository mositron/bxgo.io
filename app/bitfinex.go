package main

import (
	"fmt"
	"log"

	bitfinex "github.com/bitfinexcom/bitfinex-api-go/v1"
)

var bitfinex_loaded bool

func bitfinex_load() {
	if !bitfinex_loaded {
		go func() {
			bitfinex_start()
		}()
	}
}
func bitfinex_start() {
	if bitfinex_loaded {
		return
	}
	bitfinex_loaded = true
	c := bitfinex.NewClient()
	c.WebSocketTLSSkipVerify = false

	err := c.WebSocket.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}
	defer c.WebSocket.Close()

	//	book_btcusd_chan := make(chan []float64)
	//	book_ltcusd_chan := make(chan []float64)
	//trades_chan := make(chan []float64)
	btc := make(chan []float64)
	das := make(chan []float64)
	eth := make(chan []float64)
	omg := make(chan []float64)
	xrp := make(chan []float64)
	c.WebSocket.AddSubscribe(bitfinex.ChanTicker, bitfinex.BTCUSD, btc)
	c.WebSocket.AddSubscribe(bitfinex.ChanTicker, "DSHUSD", das)
	c.WebSocket.AddSubscribe(bitfinex.ChanTicker, bitfinex.ETHUSD, eth)
	c.WebSocket.AddSubscribe(bitfinex.ChanTicker, "OMGUSD", omg)
	c.WebSocket.AddSubscribe(bitfinex.ChanTicker, bitfinex.XRPUSD, xrp)

	go bitfinex_listen(btc, "BTC")
	go bitfinex_listen(das, "DAS")
	go bitfinex_listen(eth, "ETH")
	go bitfinex_listen(omg, "OMG")
	go bitfinex_listen(xrp, "XRP")

	err = c.WebSocket.Subscribe()
	if err != nil {
		fmt.Println("bitfinex-error: ", err.Error())
		bitfinex_loaded = false
	}
}
func bitfinex_listen(in chan []float64, crn string) {
	for {
		msg := <-in
		if len(msg) == 10 {
			Bitfinex[crn] = msg
		}
	}
}
