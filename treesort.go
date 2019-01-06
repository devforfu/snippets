package main

import (
    "fmt"
    "strconv"
    "strings"
)

func main() {
    tree := NewTree([]int{5, 10, 7, 8, 1, 12, 9, 2, 1, 0, 3})
    println(tree.String())
}

type Tree struct {
    Value int
    Left, Right *Tree
}

func NewTree(arr []int) Tree {
    root := Tree{Value:arr[0]}
    for _, x := range arr[1:] {
        root.Insert(x)
    }
    return root
}

func (n *Tree) Insert(x int) {
    if x > n.Value {
        if n.Left == nil {
            n.Left = &Tree{Value:x}
        } else {
            n.Left.Insert(x)
        }
    } else {
        if n.Right == nil {
            n.Right = &Tree{Value:x}
        } else {
            n.Right.Insert(x)
        }
    }
}

func (n *Tree) String() string {
    values := make([]string, 0)
    for _, x := range n.Sort() {
        values = append(values, strconv.Itoa(x))
    }
    return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

func (n *Tree) VisitInOrder(f func(n *Tree)) {
    if n.Left != nil { n.Left.VisitInOrder(f) }
    f(n)
    if n.Right != nil { n.Right.VisitInOrder(f) }
}

func (n *Tree) Sort() []int {
    ordered := make([]int, 0)
    n.VisitInOrder(func(n *Tree) {
        ordered = append(ordered, n.Value)
    })
    return ordered
}
