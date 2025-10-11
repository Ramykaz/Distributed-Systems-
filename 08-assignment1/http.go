package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
)

// Record represents a single log entry
type Record struct {
    Offset int    `json:"offset"`
    Value  string `json:"value"`
}

// CommitLog stores all records safely
type CommitLog struct {
    mu      sync.Mutex
    records []Record
}

// Produce adds a record
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

// List returns all records
func (c *CommitLog) List() []Record {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.records
}

// Clear removes all records
func (c *CommitLog) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.records = []Record{}
}

func main() {
    log := &CommitLog{}

    // POST / -> produce a record
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
            return
        }

        var req struct {
            Record struct {
                Value string `json:"value"`
            } `json:"record"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        record := log.Produce(req.Record.Value)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(record)
    })

    // Combined GET /records and DELETE /records
    http.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            records := log.List()
            resp := struct {
                Records []Record `json:"records"`
            }{Records: records}

            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(resp)

        case http.MethodDelete:
            log.Clear()
            w.WriteHeader(http.StatusOK)

        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    fmt.Println("Assignment 1 commit log server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
