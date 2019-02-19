package main

import "fmt"

func main() {
    naturals := make(chan int)
    processed := make(chan int)

    go counter(1,10, naturals)
    go processor(func(x int) int {
        return x * x
    }, processed, naturals)
    printer(processed)
}

func counter(start, end int, out chan<- int) {
    for x := start; x < end; x++ {
        out <- x
    }
    close(out)
}

func processor(callback func(x int)int, out chan<- int, in <-chan int) {
    for v := range in {
        out <- callback(v)
    }
    close(out)
}

func printer(in <-chan int) {
    for v := range in {
        fmt.Println(v)
    }
}
