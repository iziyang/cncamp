package main

import "fmt"

func main() {
	a := [5]string{"I", "am", "stupid", "and", "weak"}
	for _, value := range a {
		fmt.Println(value)
	}
	a[2] = "smart"
	a[4] = "strong"
	for _, value := range a {
		fmt.Println(value)
	}
}
