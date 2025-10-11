package main

import "fmt"

func main(){
	// if else exmaple
	age :=20
	if age <=13 {
	    fmt.Println("You are a child.")
	}else if age < 20 {
	    fmt.Println("You are a Teenager.")
	}else {
	    fmt.Println("You are an Adult.")
	}

    // for loop example
    fmt.Println("Counting from 1 to 5:")
    for i := 1; i <= 5; i++ {
        fmt.Println(i)
    }

    // switch example
    day := "Tuesday"
    switch day {
    case "Monday":
        fmt.Println("Start of the week")
    case "Tuesday", "Wednesday":
        fmt.Println("Midweek")
    default:
        fmt.Println("Another day")
    }
}