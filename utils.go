/*

utils.go -
the utility functions used in raid-utils

*/

package main

import (
	// internals
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// get input from console easily
func input(prompt string) string {

	// show the prompt
	fmt.Printf(prompt)

	// get input and return it
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()

}

// function to generate a random integer within a given range
func randint(min, max int) int {

	// create a seed that is different each time
	rand.Seed(time.Now().Unix())

	// generate a random int within the range
	return rand.Intn(max-min) + min

}
