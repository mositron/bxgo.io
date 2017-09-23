package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type MyHandler struct {
	sync.Mutex
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "BXGo v. "+VERSION)
	CN := map[string]int64{"BTC": 1, "ETH": 21, "DAS": 22, "XRP": 25, "OMG": 26}
	url := strings.Split(strings.TrimLeft(r.URL.Path, "/"), "/")
	cmd := url[0]
	if cmd == "ajax" {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		if pair, err := strconv.ParseInt(r.URL.Query().Get("pair"), 10, 64); err == nil && pair > 0 {
			if len(url) == 2 {
				ajax := make(map[string]interface{})
				if url[1] == "order" {
					success := false
					message := "รหัสสำหรับคำสั่งไม่ถูกต้อง"
					target := ".navbar-brand"
					action := r.URL.Query().Get("action")
					password := r.URL.Query().Get("pass")
					if password != "" && password == Conf.Pass {
						if action == "cancel" {
							if id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64); err == nil && id > 0 {
								for i := range Bot[pair].Order {
									if Bot[pair].Order[i].ID == id {
										target = "#order"
										success = true
										message = "ลบคำสั่งเรียบร้อยแล้ว... กรุณารอซักครู่เพื่ออัพเดทข้อมูล."
										Bot[pair].Order = append(Bot[pair].Order[:i], Bot[pair].Order[i+1:]...)
										api_cancel(pair, id)
										break
									}
								}
							}
						} else if action == "config" {
							LoadIni()
							target = "#delay"
							success = true
							message = "รีโหลดเรียบร้อยแล้ว... กรุณารอซักครู่เพื่ออัพเดทข้อมูล."
						}
					}
					ajax["success"] = success
					ajax["message"] = message
					ajax["target"] = target
				} else if url[1] == "data" {
					list := map[string]AList{}
					for i := range Bot {
						list[_is(i)] = AList{ID: Bot[i].Pair.ID,
							Primary:   Bot[i].Pair.Primary,
							Secondary: Bot[i].Pair.Secondary,
							Price:     Bot[i].Pair.Price,
							Change:    Bot[i].Pair.Change}
					}
					ajax["sort"] = Sort
					ajax["list"] = list
					ajax["pair"] = Bot[pair].Pair
					ajax["order"] = Bot[pair].Order
					ajax["wallet"] = Balance
					ajax["trend"] = Bot[pair].Trend
					ajax["conf"] = Bot[pair].Conf
					ajax["trans"] = Bot[pair].Trans
					ajax["sims"] = Bot[pair].Sims
					ajax["delay"] = Bot[pair].Delay
					ajax["usdthb"] = USDTHB
					ajax["bitfinex"] = Bitfinex
				}
				b, err := json.Marshal(ajax)
				if err != nil {
					_err("Web Handler - json.Marshal - ", err.Error(), "\n")
					fmt.Println(ajax)
					return
				}
				cb := r.URL.Query().Get("callback")
				if cb != "" {
					by := []byte{}
					by = append(by, cb...)
					by = append(by, "("...)
					by = append(by, b...)
					by = append(by, ")"...)
					w.Write(by)
				} else {
					w.Write(b)
				}
			}
		}
	} else if cmd == "theme" {
		dat, err := ioutil.ReadFile(r.URL.Path[1:])
		if err == nil {
			ext := filepath.Ext(r.URL.Path)
			exts := map[string]string{
				".css":  "text/css; charset=utf-8",
				".js":   "text/javascript; charset=utf-8",
				".png":  "image/png",
				".jpg":  "image/jpeg",
				".html": "text/html; charset=utf-8",
			}
			if tp, ok := exts[ext]; ok {
				w.Header().Add("Content-Type", tp)
				w.Write([]byte(dat))
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(r.URL.Path + "\nFile Not Found\n----------------------------\nBXGo v. " + VERSION + "\npowered by positron@jarm.com"))
		}
	} else if cmd == "screenshot.png" {
		http.Redirect(w, r, "/theme/default/img/screenshot.png", 301)
	} else if pa, ok := CN[cmd]; ok {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		dat, err := ioutil.ReadFile("theme/" + Conf.Theme + "/html/main.html")
		ct := strings.Replace(string(dat), "{VERSION}", VERSION, -1)
		ct = strings.Replace(ct, "{PAIR}", _is(pa), -1)
		if err == nil {
			w.Write([]byte(ct))
		} else {
			w.Write([]byte(err.Error()))
		}
	} else if cmd != "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(r.URL.Path + "\nFile Not Found\n----------------------------\nBXGo v. " + VERSION + "\npowered by positron@jarm.com"))
	} else {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		dat, err := ioutil.ReadFile("theme/" + Conf.Theme + "/html/main.html")
		ct := strings.Replace(string(dat), "{VERSION}", VERSION, -1)
		ct = strings.Replace(ct, "{PAIR}", _is(Sort[0]), -1)
		if err == nil {
			w.Write([]byte(ct))
		} else {
			w.Write([]byte(err.Error()))
		}
	}
}
