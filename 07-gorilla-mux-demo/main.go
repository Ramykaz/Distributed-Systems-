package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "sync"

    "github.com/gorilla/mux"
)

// Record represents a log entry
type Record struct {
    Offset int    `json:"offset"`
    Value  string `json:"value"`
}

// CommitLog with mutex for thread-safe access
type CommitLog struct {
    mu      sync.Mutex
    records []Record
}

// Produce adds a record safely
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

// Consume by offset safely
func (c *CommitLog) Consume(offset int) (Record, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if offset < 0 || offset >= len(c.records) {
        return Record{}, false
    }
    return c.records[offset], true
}

// List all records
func (c *CommitLog) List() []Record {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.records
}

func main() {
    log := &CommitLog{}
    r := mux.NewRouter()

    // POST /produce
    r.HandleFunc("/produce", func(w http.ResponseWriter, r *http.Request) {
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
    }).Methods("POST")

    // GET /consume?offset=0
    r.HandleFunc("/consume", func(w http.ResponseWriter, r *http.Request) {
        offsetStr := r.URL.Query().Get("offset")
        if offsetStr == "" {
            http.Error(w, "Offset required", http.StatusBadRequest)
            return
        }
        offset, err := strconv.Atoi(offsetStr)
        if err != nil {
            http.Error(w, "Invalid offset", http.StatusBadRequest)
            return
        }
        record, ok := log.Consume(offset)
        if !ok {
            http.Error(w, "Offset not found", http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(record)
    }).Methods("GET")

    // GET /consume/{offset} (URL parameter)
    r.HandleFunc("/consume/{offset}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        offsetStr := vars["offset"]
        offset, err := strconv.Atoi(offsetStr)
        if err != nil {
            http.Error(w, "Invalid offset", http.StatusBadRequest)
            return
        }
        record, ok := log.Consume(offset)
        if !ok {
            http.Error(w, "Offset not found", http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(record)
    }).Methods("GET")

    // GET /records
    r.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(log.List())
    }).Methods("GET")

    fmt.Println("Lesson 7 Gorilla Mux server running at http://localhost:8080")
    http.ListenAndServe(":8080", r)
}
