package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strconv"
//    "reflect"
)

type BinanceResponse struct {
    Symbol  string
    Price   string
}


func Cost(symbol string) float64 {
    binanceDomain := "https://api3.binance.com"
    binanceRequest := binanceDomain + "/api/v3/ticker/price?symbol=" + symbol
    resp, _ := http.Get(binanceRequest)
    defer resp.Body.Close()

    respBody, _ := ioutil.ReadAll(resp.Body)

    fmt.Println(string(respBody))

    r := &BinanceResponse{}
    json.Unmarshal(respBody, r)
    fmt.Printf("struct: %#v\n", r)

    price, _ := strconv.ParseFloat(r.Price, 64)

    fmt.Println(price)
    return price
    //fmt.Printf("http.Get body %#v\n\n\n", string(respBody))
}

