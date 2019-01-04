package main

import (
    "./src/toposort"
    "fmt"
)


func main() {
    //valid := toposort.Graph {
    //    "algorithms":            {"data structures"},
    //    "calculus":              {"linear algebra"},
    //    "compilers":             {"data structures", "formal languages", "computer organization"},
    //    "data structures":       {"discrete math"},
    //    "databases":             {"data structures"},
    //    "discrete math":         {"intro to programming"},
    //    "formal languages":      {"discrete math"},
    //    "networks":              {"operating systems"},
    //    "operating systems":     {"data structures", "computer organization"},
    //    "programming languages": {"data structures", "computer organization"},
    //}
    //
    //for i, course := range toposort.MustSort(valid) {
    //    fmt.Printf("%02d\t%v\n", i+1, course)
    //}

    faulty := toposort.Graph {
        "A": {"B", "C"},
        "B": {"C", "D"},
        "D": {"A"},
    }

    sorted, err := toposort.Sort(faulty)
    fmt.Printf("sorted=%v\nerr=%v", sorted, err)
}
