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

const (
	// random string generation things
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// useful for human-friendly input
var convertYNToBool = map[string]bool{"yes": true, "no": false}

// generates random alphabetical strings quickly
func randalphastring(n int) string {

	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {

		if remain == 0 {

			cache, remain = src.Int63(), letterIdxMax

		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {

			b[i] = letterBytes[idx]
			i--

		}

		cache >>= letterIdxBits
		remain--

	}

	return string(b)

}

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

	src := rand.New(rand.NewSource(time.Now().Unix()))
	return src.Intn(max-min) + min

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

// searchForOpenChannel searches the channel list for a channel that can take another webhook
func searchForOpenChannel() (int, error) {

	var (
		allChannels []*discordgo.Channel
		hooks       []*discordgo.Webhook
		err         error
	)

	allChannels, err = dg.GuildChannels(server.ID)
	if err != nil {

		return -1, err

	}

	channels = []*discordgo.Channel{}

	for _, channel := range allChannels {

		if channel.Type == discordgo.ChannelTypeGuildText {

			channels = append(channels, channel)

		}

	}

	for i, channel := range channels {

		hooks, err = dg.ChannelWebhooks(channel.ID)
		if err != nil {

			return -1, err

		}

		if len(hooks) != 10 {

			return i, nil

		}

	}

	return -1, nil

}
