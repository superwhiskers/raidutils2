/*

raid-utils.go -
the main file of the raid-utils rewrite

*/

package main

import (
	// internals
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"runtime"
	// externals
	"github.com/bwmarrin/discordgo"
)

// the server variable
var server *discordgo.UserGuild

// main func
func main() {

	runtime.GOMAXPROCS(1000)

	// ask for a token
	fmt.Printf("please input a discord token to use below\n")

	// initiate discordgo
	dg, err := discordgo.New("Bot " + input(": "))
	if err != nil {

		// error when trying to create a session
		fmt.Printf("[err]: could not initiate a discord session...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	// initiate the websocket connection
	err = dg.Open()
	if err != nil {

		// could not initate the websocket connection
		fmt.Printf("[err]: could not initiate the websocket connection...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	// get the servers
	servers, err := dg.UserGuilds(100, "", "")
	if err != nil {

		// could not retrieve the guilds
		fmt.Printf("[err]: could not retrieve the guilds for the bot...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	// print the servers
	for i, s := range servers {

		// print the data from the server
		fmt.Printf("%d: %s\n", i, s.Name)

	}

	// have the user select a server
	fmt.Printf("select a server...\n")
	serverIndStr := input(": ")

	// convert it to an int
	serverIndInt, err := strconv.Atoi(serverIndStr)
	if err != nil {

		// error while converting it to an int
		fmt.Printf("[err]: %s is not a number...\n", serverIndStr)
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	// ensure that index is in the slice
	if len(servers) > serverIndInt {

		// it is in the slice
		server = servers[serverIndInt]

	} else {

		// it is not in the slice
		fmt.Printf("[err]: %s is not in the server list...\n", serverIndStr)
		os.Exit(1)

	}

	// now we can start the main loop
	for {

		// show the menu
		fmt.Printf("------------------------------------\n")
		fmt.Printf(" raid-utils v2 | by superwhiskers   \n")
		fmt.Printf("------------------------------------\n")
		fmt.Printf("                                    \n")
		fmt.Printf(" (1) - dump audit logs              \n")
		fmt.Printf(" (2) - spam a server with a message \n")
		fmt.Printf(" (3) - clear webhooks in a server   \n")
		fmt.Printf(" (4) - exit                         \n")
		fmt.Printf("                                    \n")
		fmt.Printf("------------------------------------\n")
		option := input(": ")

		// check the options
		if option == "1" {

			// no audit log dumping yet

		} else if option == "2" {

			// confirm the user is fine with this
			code := strconv.Itoa(randint(100000, 999999))
			fmt.Printf("-- warning ------\n")
			fmt.Printf("you could potentially get banned from discord or\n")
			fmt.Printf("the server you are using this in with this feature...\n")
			fmt.Printf("enter the following code to perform this action: %s\n", code)
			confirmation := input(": ")

			if code == confirmation {

				// ask some configuration questions
				fmt.Printf("how many webhooks do you want to run?\n")
				hooksStr := input(": ")
				fmt.Printf("what should the webhooks be named?\n")
				name := input(": ")
				fmt.Printf("what should they spam?\n")
				message := input(": ")
				fmt.Printf("what should the avatar url be?\n")
				avatarUrl := input(": ")
				fmt.Printf("spawning webhooks and starting spam...\n")

				// convert the number of hooks to an int
				hooks, err := strconv.Atoi(hooksStr)
				if err != nil {

					// unable to perform the conversion
					fmt.Printf("[err]: unable to convert the number of hooks to an int...\n")
					fmt.Printf("       %v\n", err)
					os.Exit(1)

				}

				// get the channels from the server
				allChannels, err := dg.GuildChannels(server.ID)
				if err != nil {

					fmt.Printf("[err]: unable to retrieve channels for the selected server...\n")
					fmt.Printf("       %v\n", err)
					os.Exit(1)

				}

				// initiate the channel slice
				channels := []*discordgo.Channel{}

				// filter out all non-text channels
				for _, channel := range allChannels {

					// check the type of the channel
					if channel.Type == discordgo.ChannelTypeGuildText {

						// it is a text channel
						channels = append(channels, channel)

					}

				}

				// initiate the webhook goroutines
				for w := 1; w <= hooks; w++ {

					// combat them all picking the same channel
					time.Sleep(750 * time.Millisecond)

					// spawn a webhookWorker
					go webhookWorker(dg, channels, name, message, avatarUrl)

				}

				// wait until the user presses ctrl-c
				fmt.Printf("the spamming has begun. press ctrl-c to stop...\n")
				sc := make(chan os.Signal, 1)
				signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
				<-sc

				// if they do, let them know
				fmt.Printf("\nclosing discord connection and exiting...\n")

				// and close the connection
				dg.Close()

				// and then exit
				os.Exit(0)

			} else {

				fmt.Printf("your code is incorrect...\n")

			}

		} else if option == "3" {

			// get all webhooks in the server
			hooks, err := dg.GuildWebhooks(server.ID)
			if err != nil {

				// unable to get the server's webhooks
				fmt.Printf("[err]: unable to retrieve webhooks for the selected server...\n")
				fmt.Printf("       %v\n", err)
				os.Exit(1)

			} 

			// clear them all
			for _, hook := range hooks {

				// delete the webhook
				err := dg.WebhookDelete(hook.ID)
				if err != nil {

					// unable to delete webhook
					fmt.Printf("[err]: unable to delete webhook...\n")
					fmt.Printf("       %v\n", err)

				}

			}

		} else if option == "4" {

			// exit the program
			os.Exit(0)

		}
	}
}
