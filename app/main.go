package main

import (
	"log"
	"math"
	"time"
)

func main() {
	VERSION = "0.4.5"
	log.SetFlags(0)
	LoadIni()
	if api_key == "" || api_secret == "" {
		_tn("set api_key/api_secret.")
		for {
			time.Sleep(time.Second)
		}
	}
	USDTHB = AUSDTHB{}
	Bitfinex = map[string][]float64{}
	bitfinex_load()
	n := int64(1)
	for i := range Bot {
		Bot[i].Delay.Next_Buy = 60
		Bot[i].Delay.Next_Sell = 60
		Bot[i].Delay.Refresh_Book = n
		Bot[i].Delay.Refresh_Trade = n
		Bot[i].Delay.Refresh_Order = n
		n++
	}

	if api_port > 0 {
		w := Http{}
		go func() {
			w.Listen(api_port)
		}()
	}

	_tn("\n", _r("#", 53))
	_tn("BXGo v. ", VERSION)
	_tn("Time: " + time.Now().Format(time.Stamp))
	if api_port > 0 {
		_tn("Web Listen: ", _is(api_port))
	}
	_tn(_r("#", 53))

	for {
		use := int64(0)
		ds := ""
		if Delay.Refresh_Pair == 0 { // global
			go func() {
				api_pair()
			}()
			ds += "pair:"
			use = use_time(use, 100)
		}
		if Delay.Refresh_Balance == 0 {
			go func() {
				api_balance()
			}()
			ds += "balance:"
			use = use_time(use, 100)
		}
		if Delay.Refresh_History == 0 {
			go func() {
				api_history()
			}()
			ds += "history:"
			use = use_time(use, 100)
		}
		if Delay.Refresh_USDTHB == 0 {
			go func() {
				api_usdthb()
			}()
			ds += "usdthb:"
			use = use_time(use, 100)
		}

		for pair := range Bot {
			if Bot[pair].Delay.Refresh_Book == 0 {
				go func() {
					api_book(pair)
				}()
				ds += _is(pair) + "-book:"
				use = use_time(use, 100)
			}
			if Bot[pair].Delay.Refresh_Trade == 0 {
				go func() {
					api_trade(pair)
				}()
				ds += _is(pair) + "-trade:"
				use = use_time(use, 100)
			}
			if Bot[pair].Delay.Refresh_Order == 0 {
				go func() {
					api_order(pair)
				}()
				ds += _is(pair) + "-order:"
				use = use_time(use, 100)
			}
			process(pair)

			if Bot[pair].Delay.Next_Buy > 0 {
				Bot[pair].Delay.Next_Buy--
			}
			if Bot[pair].Delay.Next_Sell > 0 {
				Bot[pair].Delay.Next_Sell--
			}
			if Bot[pair].Delay.Refresh_Book > 0 {
				Bot[pair].Delay.Refresh_Book--
			}
			if Bot[pair].Delay.Refresh_Trade > 0 {
				Bot[pair].Delay.Refresh_Trade--
			}
			if Bot[pair].Delay.Refresh_Order > 0 {
				Bot[pair].Delay.Refresh_Order--
			}
		}
		if Delay.Refresh_Balance > 0 {
			Delay.Refresh_Balance--
		}
		if Delay.Refresh_History > 0 {
			Delay.Refresh_History--
		}
		if Delay.Refresh_Pair > 0 {
			Delay.Refresh_Pair--
		}
		iu := use
		if use > 1000 {
			use = 100
		} else {
			use = 1000 - use
		}
		if iu != 0 {
			//_tn(time.Now().Format(time.Stamp), " : ", ds, "use = ", _is(iu), " - delay = ", _is(use))
		}
		time.Sleep(time.Duration(use) * time.Millisecond)
		bitfinex_load()
	}
}

func use_time(cur int64, i int64) int64 {
	time.Sleep(time.Duration(i) * time.Millisecond)
	return cur + i
}

func process(pair int64) {
	Bot[pair].Sims = []ASims{}
	p := 0.0
	if Bot[pair].Min_Sell > 0 {
		pc := (100 / (Bot[pair].Conf.Cycle + 100)) * Bot[pair].Min_Sell
		p = ((100 / (Bot[pair].Conf.Margin + 100)) * pc)
	} else {
		p = ((100 / ((Bot[pair].Conf.Margin / 2) + 100)) * Bot[pair].Pair.Price)
	}

	if p > Bot[pair].Pair.Price {
		p = Bot[pair].Pair.Price
	}
	if p > Bot[pair].Conf.Max_Price {
		p = Bot[pair].Conf.Max_Price
	}
	for i := 0; i < 3; i++ {
		sim := _calc(pair, p)
		if sim.Buy > 0 && sim.Sell > 0 {
			Bot[pair].Sims = append(Bot[pair].Sims, sim)
			p = sim.Buy - sim.Diff
			if i == 0 && Bot[pair].Conf.Enable {
				if Bot[pair].Delay.Next_Buy == 0 && Bot[pair].Conf.Budget > 0 && sim.Buy > 0 && (Bot[pair].Conf.Max_Price >= sim.Buy || Bot[pair].Conf.Max_Price == 0) && (sim.Buy+sim.Diff > Bot[pair].Pair.Price) {
					if int64(len(Bot[pair].Order)) < Bot[pair].Conf.Max_Order {
						keep := 0.0
						/*
							if v, ok := Conf.KEEP[Pair.Primary]; ok {
								keep = v
							}
						*/
						if G_Balance[Bot[pair].Pair.Primary].Available-Bot[pair].Conf.Budget >= keep {
							if Bot[pair].Conf.Budget <= G_Balance[Bot[pair].Pair.Primary].Available {
								if sim.Order_Buy == 0 && sim.Order_Sell == 0 {
									_tn(time.Now().Format(time.Stamp)+" : send buy - ", _fs(Bot[pair].Conf.Budget), " <= ", _fs(G_Balance[Bot[pair].Pair.Primary].Available), " - rate: ", _fs(_price(pair, sim.Buy)))
									_tn("price = ", _fs(Bot[pair].Pair.Price), " , ")
									api_buy(pair, Bot[pair].Conf.Budget, _price(pair, sim.Buy))
								}
							}
						}
					}
				}
			}
		}
	}

	if Bot[pair].Conf.Enable && Bot[pair].Delay.Next_Sell == 0 && Bot[pair].Pair.Price > 0 && int64(len(Bot[pair].Order)) < Bot[pair].Conf.Max_Order && G_Balance[Bot[pair].Pair.Secondary].Available > 0 {
		keep := 0.0
		/*
			if v, ok := Conf.KEEP[Pair.Secondary]; ok {
				keep = v
			}
		*/
		if G_Balance[Bot[pair].Pair.Secondary].Available > keep {
			sell := G_Balance[Bot[pair].Pair.Secondary].Available - keep
			if sell > 0 {
				sim := _calc(pair, Bot[pair].Pair.Price)
				rate := sim.Sell
				if rate < Bot[pair].Pair.Price {
					rate = Bot[pair].Pair.Price
				}
				_tn(time.Now().Format(time.Stamp)+" : send sell - ", _fs(sell), " - rate: ", _fs(_price(pair, rate)))
				_tn("price = ", _fs(Bot[pair].Pair.Price), " , ")
				api_sell(pair, sell, _price(pair, rate))
			}
		}
	}

}

func _near(pair int64, p float64, d float64, ty string) float64 {
	for i := range Bot[pair].Order {
		if Bot[pair].Order[i].Type == ty {
			if math.Abs(Bot[pair].Order[i].Rate-p) < d {
				return Bot[pair].Order[i].Rate
			}
		}
	}
	return 0.0
}

func _calc(pair int64, p float64) ASims {
	sim := ASims{}
	sim.Buy = p
	sim.Margin = (p * (Bot[pair].Conf.Margin / 100))
	sim.Sell = sim.Margin + p
	sim.Coin = Bot[pair].Conf.Budget / p
	sim.Profit = sim.Coin * sim.Margin
	sim.Diff = (p * (Bot[pair].Conf.Cycle / 100))
	sim.Order_Sell = _near(pair, sim.Sell, sim.Diff, "sell")
	sim.Order_Buy = _near(pair, sim.Buy, sim.Diff, "buy")
	return sim
}
