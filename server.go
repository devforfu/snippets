package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "sync"
)

func main() {
    db := store{"shoes": 50, "socks": 5}
    mux := http.NewServeMux()
    mux.Handle("/", http.HandlerFunc(notFound))
    mux.Handle("/list", http.HandlerFunc(db.list))
    mux.Handle("/price", http.HandlerFunc(db.price))
    mux.Handle("/update", http.HandlerFunc(db.update))
    log.Fatal(http.ListenAndServe("localhost:8080", mux))
}

type dollars int64
type store map[string]dollars
type form map[string]string
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
    encoder := json.NewEncoder(w)
    if !ok {
        w.WriteHeader(http.StatusNotFound)
        msg := response{"error": fmt.Sprintf("no such item: %q\n", item)}
        encoder.Encode(&msg)
    } else {
        encoder.Encode(response{"success": true, "item": item, "price": price})
    }
}

func (db store) update(w http.ResponseWriter, req *http.Request) {
    var data form
    json.NewDecoder(req.Body).Decode(&data)
    encoder := json.NewEncoder(w)
    err := checkParameters(data, "item", "price")
    if err != nil {
        sendError(w, encoder, err)
        return
    }
    key := data["item"]
    price := data["price"]
    value, err := strconv.ParseInt(price, 10, 64)
    if err != nil {
        sendError(w, encoder, err)
        return
    }
    var mux sync.Mutex
    oldPrice := db[key]
    mux.Lock()
    db[key] = dollars(value)
    mux.Unlock()
    log.Printf("Price for item '%s' updated from %v to %v", key, oldPrice, value)
    encoder.Encode(response{"item": key, "old": oldPrice, "new": value})
}

func checkParameters(data form, keys ...string) error {
    for _, key := range keys {
        if _, ok := data[key]; !ok {
            return fmt.Errorf("key is missing: %s", key)
        }
    }
    return nil
}

func sendError(w http.ResponseWriter, e *json.Encoder, err error) {
    w.WriteHeader(http.StatusNotFound)
    msg := response{"error": err.Error()}
    e.Encode(&msg)
}

func notFound(w http.ResponseWriter, req *http.Request) {
    log.Printf("Not found page request: [%s] %s", req.RemoteAddr, req.RequestURI)
    json.NewEncoder(w).Encode(response{"error": "not found"})
}
