package main

import "fmt"

func hybridMode() {

	var addr string
	
	fmt.Printf("sadily, this does not work yet :p\n")

	if question("are you sure?", []string{"yes", "no"}) == "no" {

		return

	}

	for {

		addr = question("enter the address of the master", []string{})

		// attempt to verify the connection
		fmt.Printf("address: %s\n", addr)

		return

	}

}
