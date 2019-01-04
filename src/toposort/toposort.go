// Topological sorting
package toposort

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
    for i, course := range Sort(courses) {
        fmt.Printf("%02d\t%v\n", i+1, course)
    }
}

type Vertex struct {
    Item string
    ChildrenCount int
}

func (v Vertex) Order() int { return v.ChildrenCount }

// Sort topologically sorts graph.
// The vertices without direct dependencies are printed in alphabetical order.
func Sort(m map[string][]string) []string {
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
    orderedGroups := make(map[int][]string)
    maxOrder := 0

    for _, item := range order {
        vertex := item.(Vertex)
        formatted := fmt.Sprintf("%s (%d)", vertex.Item, vertex.ChildrenCount)
        group := orderedGroups[vertex.Order()]
        orderedGroups[vertex.Order()] = append(group, formatted)
        maxOrder = maxInt(vertex.Order(), maxOrder)
    }

    orderedStrings := make([]string, 0)
    for i := 0; i <= maxOrder; i++ {
        group := orderedGroups[i]
        sort.Strings(group)
        orderedStrings = append(orderedStrings, group...)
    }

    return orderedStrings
}

func maxInt(a, b int) int {
    if a >= b { return a }
    return b
}