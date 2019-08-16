package main

import (
	"fmt"
	hmw "hemmingway/io"
)


type POA struct {
	Postcode int
	State string
	Population int
}



func main() {
	results := hmw.ReadCSV("poa_population.csv")


	fmt.Print(results)


}
