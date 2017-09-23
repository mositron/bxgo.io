package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func LoadIni() {
	file, err := os.Open("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	Bot = map[int64]*ABot{}
	Sort = []int64{}
	Conf = GConf{}
	//Conf.KEEP = make(map[string]float64)
	//	coin := []string{"BCH", "BTC", "DAS", "DOG", "ETH", "FTC", "GNO", "HYP", "LTC", "NMC", "OMG", "PND", "PPC", "QRK", "REP", "THB", "XCN", "XPM", "XPY", "XRP", "ZEC"}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := strings.SplitN(strings.TrimSpace(scanner.Text()), "=", 2)
		if len(txt) == 2 {
			cmd := strings.TrimSpace(txt[0])
			val := strings.TrimSpace(txt[1])
			pair := strings.SplitN(cmd, "_", 2)
			if strings.ToUpper(cmd) == cmd {
				//Conf.KEEP[cmd], _ = strconv.ParseFloat(val, 64)
			} else if len(pair) == 2 {
				if p, err := strconv.ParseInt(pair[0], 10, 64); err == nil && p > 0 {
					if _, ok := Bot[p]; !ok {
						Bot[p] = &ABot{}
						Sort = append(Sort, p)
					}
					//Bot[p]
					cmd = strings.TrimSpace(pair[1])
					switch cmd {
					case "enable":
						Bot[p].Conf.Enable, _ = strconv.ParseBool(val)
					case "budget":
						Bot[p].Conf.Budget, _ = strconv.ParseFloat(val, 64)
					case "max_price":
						Bot[p].Conf.Max_Price, _ = strconv.ParseFloat(val, 64)
					case "max_order":
						Bot[p].Conf.Max_Order, _ = strconv.ParseInt(val, 10, 64)
					case "margin":
						Bot[p].Conf.Margin, _ = strconv.ParseFloat(val, 64)
					case "cycle":
						Bot[p].Conf.Cycle, _ = strconv.ParseFloat(val, 64)
					}
				}
			} else {
				switch cmd {
				case "port":
					api_port, _ = strconv.ParseInt(val, 10, 64)
				case "key":
					api_key = val
				case "secret":
					api_secret = val
				case "pass":
					api_pass = val
				case "theme":
					Conf.Theme = val
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
