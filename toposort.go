// Topological sorting
package main

import "fmt"

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

func (g *Graph) TopologicalSort() []string {
    seen := make(map[string]bool)
    order := make([]string, 0)
    var visit func (string)

    visit = func (item string) {

    }
}

func main() {
    for i, course := range TopologicalSort(courses) {
        fmt.Printf("%02d: %v\n", i+1, course)
    }
}

type Vertex struct {
    Item string
    ChildrenCount int
}

type Vertices = []Vertex

func (vs *Vertices) Len() int { return len(*vs) }

func (vs *Vertices) Less(i, j int) bool { return (*vs)[i].ChildrenCount < (*vs)[j].ChildrenCount }

func (vs *Vertices) Swap(i, j int) { (*vs)[i], (*vs)[j] = (*vs)[j], (*vs)[i] }

func TopologicalSort(m map[string][]string) [] string {
    seen := make(map[string]bool)
    order := make([]Vertex, 0)
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


    return order
}

