package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

type Client struct {
    baseURL string
}

type JSON map[string]interface{}
type Price float64
type PriceList map[string]Price

func (c *Client) Parse() {
    host := flag.String("-h", "localhost", "host address")
    port := flag.Int("-p", 8080, "port to bind")
    flag.Parse()
    c.baseURL = fmt.Sprintf("http://%s:%d", *host, *port)
}

func (c *Client) FetchPriceList() (prices PriceList, err error) {
    resp, err := c.get("/list")
    prices = make(PriceList)
    for k, v := range resp { prices[k] = Price(v.(float64)) }
    return prices, err
}

func (c *Client) MustFetchPriceList() (prices PriceList) {
    prices, err := c.FetchPriceList()
    if err != nil { log.Fatalf("cannot get prices: %v", err)}
    return prices
}

func (c *Client) PrintPrices(w io.Writer, prices PriceList) {
    w.Write([]byte("List of prices:\n"))
    for k, v := range prices {
        w.Write([]byte(fmt.Sprintf("%s\t%v\n", k, v)))
    }
}

func (c *Client) get(endpoint string) (result JSON, err error) {
    url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
    resp, err := http.Get(url)
    if err != nil { return }
    defer resp.Body.Close()
    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(&result)
    return result, err
}

func main() {
    client := Client{}
    client.Parse()
    prices := client.MustFetchPriceList()
    client.PrintPrices(os.Stdout, prices)
}
