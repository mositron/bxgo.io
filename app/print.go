package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func _r(txt string, count int) string {
	return (strings.Repeat(txt, count))
}

func _p(txt string, count int, end string) {
	if end != "\n" {
		end += " "
	}
	_t((strings.Repeat(txt, count)) + " " + end)
}

func _f(i float64, count int) float64 {
	s, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(count)+"f", i), 64)
	return s
}

func _fs(i float64) string {
	return strconv.FormatFloat(_f(i, 8), 'f', -1, 64)
}

func _is(i int64) string {
	return strconv.FormatInt(i, 10)
}
func _ir(i int64, r int64) int64 {
	return (i - r) + rand.Int63n((i*2)+1)
}

func _t(s ...string) {
	fmt.Print(strings.Join(s, ""))
}

func _tn(s ...string) {
	_t(strings.Join(s, ""), "\n")
}

func _err(s ...string) {
	_tn("ERROR - Time:", time.Now().Format(time.Stamp))
	_tn("ERROR - Message:", strings.Join(s, ""), "\n")
}

func _price(pair int64, i float64) float64 {
	if Bot[pair].Pair.Primary == "THB" {
		return _f(i, 2)
	}
	return _f(i, 8)
}
