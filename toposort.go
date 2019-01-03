// Topological sorting
package main

import (
    "fmt"
    "sort"
)

type Graph = map[string][]string

var courses = Graph {
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
    for i, course := range TopologicalSort(courses) {
        fmt.Printf("%02d\t%v\n", i+1, course)
    }
}

type OrderedItem interface {
    Order() int
}

type OrderedCollection []OrderedItem
func (c OrderedCollection) Len() int { return len(c) }
func (c OrderedCollection) Less(i, j int) bool { return c[i].Order() < c[j].Order() }
func (c OrderedCollection) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

type Vertex struct {
    Item string
    ChildrenCount int
}

func (v Vertex) Order() int { return v.ChildrenCount }

func TopologicalSort(m map[string][]string) []string {
    seen := make(map[string]bool)
    order := make(OrderedCollection, 0)
    var visit func (string)

    visit = func (item string) {
        if !seen[item] {
            seen[item] = true
            for _, child := range m[item] {
                visit(child)
            }
            order = append(order, Vertex{item, len(m[item])})
        }
    }

    for key := range m { visit(key) }

    sort.Sort(order)
    orderedStrings := make([]string, 0)
    for _, item := range order {
        vertex := item.(Vertex)
        formatted := fmt.Sprintf("%s (%d)", vertex.Item, vertex.ChildrenCount)
        orderedStrings = append(orderedStrings, formatted)
    }
    return orderedStrings
}

