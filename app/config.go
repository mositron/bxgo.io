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
	Conf.URL = "https://bx.in.th/api/"

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
					cmd = strings.TrimSpace(pair[1])
					switch cmd {
					case "enable":
						if v, err := strconv.ParseBool(val); err == nil {
							Bot[p].Conf.Enable = v
						}
					case "budget":
						if v, err := strconv.ParseFloat(val, 64); err == nil && v > 0 {
							Bot[p].Conf.Budget = v
						}
					case "max_price":
						if v, err := strconv.ParseFloat(val, 64); err == nil && v > 0 {
							Bot[p].Conf.Max_Price = v
						}
					case "max_order":
						if v, err := strconv.ParseInt(val, 10, 64); err == nil && v > 0 {
							Bot[p].Conf.Max_Order = v
						}
					case "margin":
						if v, err := strconv.ParseFloat(val, 64); err == nil && v > 0 {
							Bot[p].Conf.Margin = v
						}
					case "cycle":
						if v, err := strconv.ParseFloat(val, 64); err == nil && v > 0 {
							Bot[p].Conf.Cycle = v
						}
					}
				}
			} else {
				switch cmd {
				case "port":
					if v, err := strconv.ParseInt(val, 10, 64); err == nil && v > 0 {
						Conf.Port = v
					}
				case "key":
					if val != "" {
						Conf.Key = val
					}
				case "secret":
					if val != "" {
						Conf.Secret = val
					}
				case "pass":
					if val != "" {
						Conf.Pass = val
					}
				case "theme":
					if val != "" {
						Conf.Theme = val
					}
				case "line":
					if val != "" {
						Conf.Line = val
					}
				case "twofa":
					if val != "" {
						Conf.TwoFA = val
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
