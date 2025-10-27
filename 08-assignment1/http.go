package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Record struct {
	Value  string `json:"value"`
	Offset int    `json:"offset"`
}

type CommitLog struct {
	mu      sync.RWMutex
	records []Record
}

func (c *CommitLog) Append(record Record) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record.Offset = len(c.records)
	c.records = append(c.records, record)
	return record.Offset, nil
}

func (c *CommitLog) Read(offset int) (Record, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if offset < 0 || offset >= len(c.records) {
		return Record{}, http.ErrNoLocation
	}
	return c.records[offset], nil
}

func (c *CommitLog) List() []Record {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]Record{}, c.records...) // Return copy
}

func (c *CommitLog) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.records = []Record{}
}

var commitLog = &CommitLog{}

func handleProduce(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Record Record `json:"record"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	offset, _ := commitLog.Append(req.Record)
	res := struct {
		Offset int `json:"offset"`
	}{Offset: offset}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func handleConsume(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}
	record, err := commitLog.Read(offset)
	if err != nil {
		http.Error(w, "Offset out of range", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func handleList(w http.ResponseWriter, r *http.Request) {
	records := commitLog.List()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]Record{"records": records})
}

func handleClear(w http.ResponseWriter, r *http.Request) {
	commitLog.Clear()
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleProduce(w, r)
		case http.MethodGet:
			handleConsume(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleList(w, r)
		case http.MethodDelete:
			handleClear(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
