package toposort

// OrderedItem has an integer key assigned to it and can be
// sorted using this key.
type OrderedItem interface {
    Order() int
}

// OrderedCollection is an array of ordered items conforming
// the sorting interface.
type OrderedCollection []OrderedItem
func (c OrderedCollection) Len() int { return len(c) }
func (c OrderedCollection) Less(i, j int) bool { return c[i].Order() < c[j].Order() }
func (c OrderedCollection) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
