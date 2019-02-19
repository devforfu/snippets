package main

import (
    "fmt"
    "time"
)

func main() {
    naturals := make(chan int)
    squares := make(chan int)

    go func() {
        defer close(naturals)
        for x := 1; x < 10; x++ {
            naturals <- x
        }
    }()

    go func() {
        defer close(squares)
        for {
            x, ok := <- naturals
            if !ok { break }
            squares <- x * x
        }
    }()

    for x := range squares {
        fmt.Println(x)
        time.Sleep(100*time.Millisecond)
    }
}
