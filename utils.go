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
	"strings"
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

// function that asks a question until it gets a valid answer
func question(prompt string, valid []string) string {

	var inp string

	for {

		fmt.Printf("%s\n", prompt)
		if len(valid) != 0 {

			fmt.Printf("(%s): ", strings.Join(valid, ", "))

		} else {

			fmt.Printf(": ")

		}

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		inp = scanner.Text()

		if len(valid) == 0 {

			return inp

		}

		for _, ele := range valid {

			if ele == inp {

				return inp

			}

		}

		fmt.Printf("%s is not a valid answer\n", inp)

	}

}
