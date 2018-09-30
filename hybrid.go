package main

import "fmt"

func hybridMode() {

	var addr string

	if question("are you sure?", []string{"yes", "no"}) == "no" {

		return

	}

	for {

		addr = question("enter the address of the master", []string{})

		// attempt to verify the connection

	}

}
