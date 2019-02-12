package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
)

func main() {
    host := flag.String("-h", "localhost", "host address")
    port := flag.Int("-p", 8080, "port to bind")
    flag.Parse()
    url := fmt.Sprintf("http://%s:%d", *host, *port)
    resp, err := http.Get(fmt.Sprintf("%s/list", url))
    if err != nil { log.Fatalf("cannot query the server %s: %s", url, err) }
    defer resp.Body.Close()
    decoder := json.NewDecoder(resp.Body)
    var content = make(map[string]int)
    err = decoder.Decode(&content)
    if err != nil { log.Fatalf("cannot decode server's response: %s", err) }
    fmt.Printf("Decoded JSON: %v\n", content)
}

func requestListOfItems() {

}