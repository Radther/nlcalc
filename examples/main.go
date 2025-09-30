package main

import (
	"fmt"
	"log"

	"github.com/radther/nlcalc/pkg/nlcalc"
)

func main() {
	result, err := nlcalc.Parse("ten plus fifteen", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %g\n", result)
	// Output: Result: 25
}
