package main

import (
    "./src/tempconv"
    "flag"
    "fmt"
)

var temp = tempconv.CelsiusFlag("temp", 20.0, "temperature")

func main() {
    flag.Parse()
    fmt.Println(*temp)
}
