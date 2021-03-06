package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "strings"
    "time"
)

func main()  {
    go func() {
        listener, err := net.Listen("tcp", ":8080")
        if err != nil {
            log.Fatal(err)
        }
        for {
            conn, err := listener.Accept()
            if err != nil {
                log.Print(err)
                continue
            }
            go handle(conn)
        }
    }()

    conn, err := net.Dial("tcp", ":8080")
    log.Printf("Connection type: %T", conn)
    if err != nil { log.Fatal(err) }
    done := make(chan struct{})
    go func() {
        io.Copy(os.Stdout, conn)
        log.Println("done")
        done <- struct{}{}
    }()
    mustCopy(conn, os.Stdin)
    conn.Close()
    <- done
}

func handle(c net.Conn) {
    defer c.Close()
    log.Printf("Accepted connection: %s", c.RemoteAddr().String())
    input := bufio.NewScanner(c)
    for input.Scan() {
        go echo(c, input.Text(), 1*time.Second)
    }
}

func echo(c net.Conn, shout string, delay time.Duration) {
    fmt.Fprintln(c, "\t", strings.ToUpper(shout))
    time.Sleep(delay)
    fmt.Fprintln(c, "\t", shout)
    time.Sleep(delay)
    fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func mustCopy(dst io.Writer, src io.Reader) {
    if _, err := io.Copy(dst, src); err != nil {
        log.Fatal(err)
    }
}