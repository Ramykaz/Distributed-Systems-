package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"   // provide synchronization primitives like mutex it prevents corruption of data
)

// Record represents a single log entry
type Record struct {
    Offset int    `json:"offset"`
    Value  string `json:"value"`
}

// CommitLog stores all records
type CommitLog struct {
    mu      sync.Mutex
    records []Record
}

// Produce adds a record to the log
func (c *CommitLog) Produce(value string) Record {
    c.mu.Lock()
    defer c.mu.Unlock()
    record := Record{
        Offset: len(c.records),
        Value:  value,
    }
    c.records = append(c.records, record)
    return record
}

// Consume returns a record by offset
func (c *CommitLog) Consume(offset int) (Record, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if offset < 0 || offset >= len(c.records) {
        return Record{}, false
    }
    return c.records[offset], true
}

func main() {
    log := &CommitLog{}

    // Produce endpoint
    http.HandleFunc("/produce", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
            return
        }
        var req struct {
            Value string `json:"value"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        record := log.Produce(req.Value)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(record)
    })

    // Consume endpoint
    http.HandleFunc("/consume", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
            return
        }
        var req struct {
            Offset int `json:"offset"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        record, ok := log.Consume(req.Offset)
        if !ok {
            http.Error(w, "Offset not found", http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(record)
    })

    fmt.Println("Commit log server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
