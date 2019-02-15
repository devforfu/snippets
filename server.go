package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

func main() {
    db := store{"shoes": 50, "socks": 5}
    mux := http.NewServeMux()
    mux.Handle("/", http.HandlerFunc(notFound))
    mux.Handle("/list", http.HandlerFunc(db.list))
    mux.Handle("/price", http.HandlerFunc(db.price))
    log.Fatal(http.ListenAndServe("localhost:8080", mux))
}

type dollars int64
type store map[string]dollars
type response map[string]interface{}

func (db store) list(w http.ResponseWriter, req *http.Request) {
    var prices = make(map[string]int)
    for item, price := range db {
        prices[item] = int(price)
    }
    encoder := json.NewEncoder(w)
    encoder.Encode(&prices)
}

func (db store) price(w http.ResponseWriter, req *http.Request) {
    item := req.URL.Query().Get("item")
    price, ok := db[item]
    encoder := json.NewEncoder(w);
    if !ok {
        w.WriteHeader(http.StatusNotFound)
        msg := response{"error": fmt.Sprintf("no such item: %q\n", item)}
        encoder.Encode(&msg)
    } else {
        encoder.Encode(response{"success": true, "item": item, "price": price})
    }
}

func notFound(w http.ResponseWriter, req *http.Request) {
    log.Printf("Not found page request: [%s] %s", req.RemoteAddr, req.RequestURI)
    json.NewEncoder(w).Encode(response{"error": "not found"})
}
