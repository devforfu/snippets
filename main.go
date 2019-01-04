package main

import (
    "./src/toposort"
    "fmt"
)

var graph = toposort.Graph {
    "algorithms":            {"data structures"},
    "calculus":              {"linear algebra"},
    "compilers":             {"data structures", "formal languages", "computer organization"},
    "data structures":       {"discrete math"},
    "databases":             {"data structures"},
    "discrete math":         {"intro to programming"},
    "formal languages":      {"discrete math"},
    "networks":              {"operating systems"},
    "operating systems":     {"data structures", "computer organization"},
    "programming languages": {"data structures", "computer organization"},
}


func main() {
    //var graph = toposort.Graph {
    //    "A": {"B", "C"},
    //    "B": {"C"},
    //    "C": {"Z", "X", "Y"},
    //}
    for i, course := range toposort.Sort(graph) {
        fmt.Printf("%02d\t%v\n", i+1, course)
    }
}
