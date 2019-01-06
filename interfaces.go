package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

type ByteCounter int
func (c *ByteCounter) Write(p []byte) (int, error) {
    *c += ByteCounter(len(p))
    return len(p), nil
}

type WordCounter int
func (c *WordCounter) Write(p []byte) (int, error) {
    var nWords, nBytes int
    for {
        advance, _, err := bufio.ScanWords(p, true)
        if err != nil { return 0, err }
        if advance == 0 { break }
        p = p[advance:]
        nBytes += advance - 1
        nWords++
    }
    *c += WordCounter(nWords)
    return nBytes, nil
}
func (c *WordCounter) Reset() {
    *c = 0
}

type wrapper struct {
    wrapped io.Writer
    counter int64
}
func (w *wrapper) Write(p []byte) (int, error) {
    n, err := w.wrapped.Write(p)
    if err != nil { return 0, err }
    w.counter += int64(n)
    return n, nil
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
    wrapper := wrapper{wrapped: w}
    var writer io.Writer = &wrapper
    return writer, &wrapper.counter
}

func main() {
    var c ByteCounter
    c.Write([]byte("Hello"))
    fmt.Println(c)

    c = 0
    var name = "Dolly"
    fmt.Fprintf(&c, "hello, %s", name)
    fmt.Println(c)

    var wc WordCounter
    fmt.Printf("%-30s count\n", "string")
    for _, word := range []string{
        "",
        "Single",
        "Phrase with four words",
        "StringWithoutSpaces!",
    }{
        fmt.Fprintf(&wc, word)
        fmt.Printf("%-30s % 5d\n", word, wc)
        wc.Reset()
    }

    var writer io.Writer
    writer, err := os.Create("file.txt")
    if err != nil { os.Exit(1) }
    writer, count := CountingWriter(writer)
    lines := []string{"First", "Second", "Third"}
    for _, line := range lines {
        writer.Write([]byte(line + "\n"))
    }
    fmt.Printf("Number of bytes written: %d", *count)
}
