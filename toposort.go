// Topological sorting
package main

func main() {

}

func TopologicalSort(m map[string][]string) {
    topoSort(m)
}

func topoSort(m map[string][]string) []string {
    seen := make(map[string]bool)
    order := make([]string, 0)
    var visit func (string)

    visit = func (item string) {
        if !seen[item] {
            seen[item] = true
            for _, child := range m[item] {
                visit(child)
            }
            order = append(order, item)
        }
    }

    for key := range m { visit(key) }

    return order
}

