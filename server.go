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

type JSONResponse struct {
    http.ResponseWriter
    *json.Encoder
}

func (r JSONResponse) SendError(message string) {
    r.WriteHeader(http.StatusBadRequest)
    r.Encode(response{"error": message})
}

func (r JSONResponse) SendStore(s store) {
    r.WriteHeader(http.StatusAccepted)
    r.Encode(s)
}

func (r JSONResponse) SendPrice(item string, price dollars) {
    r.WriteHeader(http.StatusAccepted)
    r.Encode(response{"item": item, "price": price})
}

func (r JSONResponse) SendSuccess(resp response) {
    r.WriteHeader(http.StatusAccepted)
    r.Encode(resp)
}

func NewJSONResponse(w http.ResponseWriter) JSONResponse {
    resp := JSONResponse{}
    resp.ResponseWriter = w
    resp.Encoder = json.NewEncoder(w)
    return resp
}

func (db store) list(w http.ResponseWriter, req *http.Request) {
    NewJSONResponse(w).SendStore(db)
}

func (db store) price(w http.ResponseWriter, req *http.Request) {
    item := req.URL.Query().Get("item")
    resp := NewJSONResponse(w)
    if price, ok := db[item]; !ok {
        resp.SendError(fmt.Sprintf("no such item: %s", item))
    } else {
        resp.SendPrice(item, price)
    }
}

func (db store) update(w http.ResponseWriter, req *http.Request) {
    var data form
    json.NewDecoder(req.Body).Decode(&data)
    resp := NewJSONResponse(w)
    err := checkParameters(data, "item", "price")
    if err != nil {
        resp.SendError(err.Error())
        return
    }
    key, price := data["item"], data["price"]
    value, err := strconv.ParseInt(price, 10, 64)
    if err != nil {
        resp.SendError(err.Error())
        return
    }
    var mux sync.Mutex
    mux.Lock()
    oldPrice := db[key]
    mux.Unlock()
    log.Printf("Price for item '%s' updated from %v to %v", key, oldPrice, value)
    resp.SendSuccess(response{"item": key, "old": oldPrice, "new": value})
}

func checkParameters(data form, keys ...string) error {
    for _, key := range keys {
        if _, ok := data[key]; !ok {
            return fmt.Errorf("key is missing: %s", key)
        }
    }
    return nil
}

func notFound(w http.ResponseWriter, req *http.Request) {
    log.Printf("Not found page request: [%s] %s", req.RemoteAddr, req.RequestURI)
    json.NewEncoder(w).Encode(response{"error": "not found"})
}
