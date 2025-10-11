package main

import "fmt"

func main() {
    // Variables
    var name string = "Tom Hardy"
    age := 34  // short declaration
    height := 1.78

    // Constants
    const country = "UK"

    // Print values
    fmt.Println("Name:", name)
    fmt.Println("Age:", age)
    fmt.Println("Height:", height)
    fmt.Println("Country:", country)

    // Using a function
    greet(name)
}

// Function definition
func greet(person string) {
    fmt.Println("Hello,", person, "welcome to Go!")
}
