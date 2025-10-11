package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// Define a struct
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // HTTP handler
    http.HandleFunc("/person", func(w http.ResponseWriter, r *http.Request) {
        p := Person{Name: "Chewbaca", Age: 125}

        // Convert struct to JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(p)
    })

    fmt.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
