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

type Vertex struct {
    Item string
    ChildrenCount int
}

func (v Vertex) Order() int { return v.ChildrenCount }

// MustSort topologically sorts graph or panics in case if cycle discovered
func MustSort(m Graph) []string {
    sorted, err := Sort(m)
    if err != nil { panic(err) }
    return sorted
}

// Sort topologically sorts graph.
// The vertices without direct dependencies are printed in alphabetical order.
// Returns error if cycle discovered.
func Sort(m Graph) ([]string, error) {
    const (
        White = iota
        Gray
        Black
    )
    seen := make(map[string]int)
    order := make(OrderedCollection, 0)

    var visit func (string) error

    visit = func (item string) (err error) {
        if seen[item] == White {
            seen[item] = Gray
            for _, child := range m[item] {
                if seen[child] == Gray {
                    return fmt.Errorf("cycle discovered: %s <--> %s", item, child)
                }
                if err = visit(child); err != nil {
                    return err
                }
            }
            order = append(order, Vertex{item, len(m[item])})
            seen[item] = Black
        }
        return
    }

    for key := range m {
        if err := visit(key); err != nil {
            return nil, err
        }
    }

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

    return orderedStrings, nil
}

func maxInt(a, b int) int {
    if a >= b { return a }
    return b
}