package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

type Track struct {
    Title, Artist, Album string
    Year int
    Length time.Duration
}

func (t *Track) String() string {
    return fmt.Sprintf("%s - %s, %d (%v)", t.Title, t.Artist, t.Year, t.Length)
}

type songs []*Track
type database map[string]songs

func createDatabase(tracks []*Track) database {
    db := make(database)
    for _, track := range tracks {
        value, ok := db[track.Title]
        if !ok {
            value = make(songs, 0)
        }
        value = append(value, track)
        db[track.Title] = value
    }
    return db
}

func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
   for _, tracks := range db {
       for _, track := range tracks {
           fmt.Printf("%v\n", track)
       }
   }
}

func main() {
    length := func (s string) time.Duration {
        d, err := time.ParseDuration(s)
        if err != nil { panic(s) }
        return d
    }
    tracks := []*Track {
       {"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
       {"Go", "Moby", "Moby", 1992, length("3m37s")},
       {"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
       {"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
    }
    db := createDatabase(tracks)
    log.Fatal(http.ListenAndServe("localhost:8080", db))
}