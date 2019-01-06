// An example from The Go Programming Language (6.5)
package main

import (
    "bytes"
    "fmt"
)

func main() {
    var set IntSet
    for _, x := range [...]int{0, 1, 2, 64, 66, 127} {
        set.Add(x)
    }
    println(set.String())
    println(Inspect(&set))
}

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
    words []uint64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
    word, bit := x/64, uint(x%64)
    return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
    word, bit := x/64, uint(x%64)
    for word >= len(s.words) {
        s.words = append(s.words, 0)
    }
    s.words[word] |= 1<<bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
    for i, tWord := range t.words {
        if i < len(s.words) {
            s.words[i] |= tWord
        } else {
            s.words = append(s.words, tWord)
        }
    }
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
    var buf bytes.Buffer
    buf.WriteByte('{')
    for i, word := range s.words {
        if word == 0 { continue }
        for j := 0; j < 64; j++ {
            if word&(1<<uint(j)) != 0 {
                if buf.Len() > len("{") {
                    buf.WriteByte(' ')
                }
                fmt.Fprintf(&buf, "%d", 64*i + j)
            }
        }
    }
    buf.WriteByte('}')
    return buf.String()
}

// Inspect returns string with raw representation of the set.
// Useful for debugging purposes.
func Inspect(s *IntSet) string {
    var buf bytes.Buffer
    for _, word := range s.words {
        for j := 0; j < 64; j++ {
            if j != 0 && j % 4 == 0 {
                buf.WriteByte('_')
            }
            if word&(1<<uint(j)) == 0 {
                buf.WriteByte('0')
            } else {
                buf.WriteByte('1')
            }
        }
        buf.WriteByte('\n')
    }
    return buf.String()
}
