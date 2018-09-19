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
	// externals
	"github.com/bwmarrin/discordgo"
)

// useful for human-friendly input
var convertYNToBool = map[string]bool{"yes": true, "no": false}

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

		fmt.Printf("\"%s\" is not a valid answer\n", inp)

	}

}

// creates a new webhook for use
func createWebhook(channel *discordgo.Channel, name string) (*discordgo.Webhook, error) {

	var sleepTime time.Duration

	webhook, err := dg.WebhookCreate(channel.ID, name, "")
	if err != nil {

		// webhook recreation loop thing
		for tries := 0; tries <= 3; tries++ {

			webhook, err = dg.WebhookCreate(channel.ID, name, "")
			if err != nil {

				fmt.Printf("[err]: unable to recreate webhook...\n")
				fmt.Printf("       %v\n", err)

				if tries == 3 {

					return &discordgo.Webhook{}, err

				} else {

					sleepTime, err = time.ParseDuration(fmt.Sprintf("%ds", 30*tries))
					if err != nil {

						fmt.Printf("[err]: unable to parse duration...\n")
						fmt.Printf("       %v\n", err)
						return &discordgo.Webhook{}, err

					}

					time.Sleep(sleepTime)

				}
				continue

			}
			break

		}

	}

	return webhook, nil

}
