package main

import "fmt"

// Define a struct
type Person struct {
    Name string
    Age  int
}

// Method for Person
func (p Person) Greet() {
    fmt.Println("Hello, my name is", p.Name, "and I am", p.Age, "years old.")
}

func main() {
    // Create a Person instance
    p1 := Person{Name: "Ramy", Age: 25}

    // Call the method
    p1.Greet()

    // Another example
    p2 := Person{Name: "Sara", Age: 30}
    p2.Greet()
}
