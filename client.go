package main

import (
    "bytes"
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

func (obj JSON) Buffer() (buf bytes.Buffer, err error) {
    encoder := json.NewEncoder(&buf)
    err = encoder.Encode(obj)
    if err != nil { return }
    return buf, nil
}

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

func (c *Client) Update(item string, price Price) (oldPrice Price, err error) {
    data := make(JSON)
    data["item"] = item
    // without casting into string, fails to marshal
    data["price"] = fmt.Sprintf("%v", price)
    buf, err := data.Buffer()
    if err != nil { return }
    update, err := c.post("/update", &buf)
    if err != nil { return }
    return Price(update["old"].(float64)), err
}

func (c *Client) PrintPrices(w io.Writer, prices PriceList) {
    w.Write([]byte("List of prices:\n"))
    for k, v := range prices {
        w.Write([]byte(fmt.Sprintf("%s\t%v\n", k, v)))
    }
}

func (c *Client) get(endpoint string) (result JSON, err error) {
    resp, err := http.Get(c.url(endpoint))
    if err != nil { return }
    defer resp.Body.Close()
    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(&result)
    return result, err
}

func (c *Client) post(endpoint string, data io.Reader) (result JSON, err error) {
    resp, err := http.Post(c.url(endpoint), "application/json", data)
    if err != nil { return }
    defer resp.Body.Close()
    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(&result)
    return result, err
}
//
//func (c *Client) postJSON(endpoint string, data JSON) (result JSON, err error) {
//    jsonValue, _ := json.Marshal(data)
//    buf := bytes.NewBuffer(jsonValue)
//    resp, err := http.Post(c.url(endpoint), "application/json", buf)
//    if err != nil { return }
//    defer resp.Body.Close()
//    decoder := json.NewDecoder(resp.Body)
//    err = decoder.Decode(&result)
//    return result, err
//}

func (c *Client) url(endpoint string) string {
    return fmt.Sprintf("%s%s", c.baseURL, endpoint)
}

func main() {
    client := Client{}
    client.Parse()
    prices := client.MustFetchPriceList()
    client.PrintPrices(os.Stdout, prices)
    oldPrice, _ := client.Update("socks", 10)
    fmt.Printf("Old socks price was: %v\n", oldPrice)
}
