package main

import "time"

type Track struct {
    Title, Artist, Album string
    Year int
    Length time.Duration
}

type songs []*Track
type database map[string]songs

func createDatabase(tracks []*Track) database {
    db := make(database)
    for _, track := range tracks {
        value, ok := db[track.Title]
        if !ok {
            value = make(songs, 0)
            db[track.Title] = value
        }
        value = append(value, track)
    }
    return db
}

//func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//    for key, track := range db {
//        fmt
//    }
//}

func main() {
    //tracks := []*Track {
    //    {"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
    //    {"Go", "Moby", "Moby", 1992, length("3m37s")},
    //    {"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
    //    {"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
    //}
    //db := createDatabase(tracks)
    //log.Fatal(http.ListenAndServe("localhost:8080", db))
}