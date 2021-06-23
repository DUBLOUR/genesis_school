package main

import (
	//    "fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	//    "reflect"
)

//type Price float64

type BinanceResponse struct {
	Symbol string
	Price  string
}

func Cost(currency string) (float64, error) {
	marketEndpoint := "https://api3.binance.com/api/v3/ticker/price?"
	params := url.Values{}
	params.Set("symbol", currency)

	r, err := http.Get(marketEndpoint + params.Encode())
	defer r.Body.Close()
	if err != nil {
		return 0, err
	}

	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}

	respData := &BinanceResponse{}
	if err = json.Unmarshal(respBody, respData); err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(respData.Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
