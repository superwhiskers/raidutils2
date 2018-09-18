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
	"runtime"
	"strconv"
	"syscall"
	// externals
	"github.com/bwmarrin/discordgo"
)

// variables that needed to be glo
var (
	server   *discordgo.UserGuild
	channels []*discordgo.Channel
	dg       *discordgo.Session
	err      error
)

// main func
func main() {

	runtime.GOMAXPROCS(1000)

	token := question("input a discord token below", []string{})

	switch question("is this a bot?", []string{"yes", "no"}) {

	case "yes":
		dg, err = discordgo.New("Bot " + token)

	case "no":
		dg, err = discordgo.New(token)

	}

	if err != nil {

		fmt.Printf("[err]: could not initiate a discord session...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	err = dg.Open()
	if err != nil {

		fmt.Printf("[err]: could not initiate the websocket connection...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	servers, err := dg.UserGuilds(100, "", "")
	if err != nil {

		fmt.Printf("[err]: could not retrieve the guilds for the bot...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	for i, s := range servers {

		fmt.Printf("%d: %s\n", i, s.Name)

	}

	serverIndStr := question("select a server", []string{})

	serverIndInt, err := strconv.Atoi(serverIndStr)
	if err != nil {

		// error while converting it to an int
		fmt.Printf("[err]: %s is not a number...\n", serverIndStr)
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	if len(servers) > serverIndInt {

		server = servers[serverIndInt]

	} else {

		fmt.Printf("[err]: %s is not in the server list...\n", serverIndStr)
		os.Exit(1)

	}

	for {

		fmt.Println("           _     _       _   _ _     ")
		fmt.Println("          (_)   | |     | | (_) |    ")
		fmt.Println(" _ __ __ _ _  __| |_   _| |_ _| |___ ")
		fmt.Println("| '__/ _` | |/ _` | | | | __| | / __|")
		fmt.Println("| | | (_| | | (_| | |_| | |_| | \\__ \\")
		fmt.Println("|_|  \\__,_|_|\\__,_|\\__,_|\\__|_|_|___/")
		fmt.Println("by superwhiskers")
		fmt.Println()
		fmt.Println(" 1. start the raid")
		fmt.Println(" 2. open tool menu")
		fmt.Println(" 3. exit")
		fmt.Println()

		switch question("select an option", []string{"1", "2", "3"}) {

		case "1":
			code := strconv.Itoa(randint(100000, 999999))
			fmt.Printf("warning~~!\n")
			fmt.Printf("you could potentially get banned from discord or\n")
			fmt.Printf("the server you are using this in with this feature...\n")
			fmt.Printf("enter the following code to perform this action: %s\n", code)

			if code == input(": ") {
				name := question("what should the webhooks be named?", []string{})
				message := question("what should they spam?", []string{})
				avatarUrl := question("what should the avatar url be?", []string{})
				fmt.Println("calculating webhook count...")

				allChannels, err := dg.GuildChannels(server.ID)
				if err != nil {

					fmt.Printf("[err]: unable to retrieve channels for the selected server...\n")
					fmt.Printf("       %v\n", err)
					os.Exit(1)

				}

				channels = []*discordgo.Channel{}

				for _, channel := range allChannels {

					if channel.Type == discordgo.ChannelTypeGuildText {

						channels = append(channels, channel)

					}

				}

				hooks := len(channels) * 10

				fmt.Printf("approximately %d webhooks will be spawned\n", hooks)
				if question("are you sure you want to proceed?", []string{"yes", "no"}) == "no" {

					os.Exit(0)

				}

				fmt.Printf("spawned 0/%d workers...", hooks)

				for w := 1; w <= hooks; w++ {

					go webhookWorker(w-1, name, message, avatarUrl)

					fmt.Printf("\rspawned %d/%d workers...", w, hooks)

				}

				fmt.Printf("\nthe spamming has begun. press ctrl-c to stop...\n")
				sc := make(chan os.Signal, 1)
				signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
				<-sc

				fmt.Printf("\nclosing discord connection and exiting...\n")
				dg.Close()
				os.Exit(0)

			} else {

				fmt.Printf("your code is incorrect...\n")

			}

		case "2":
			fmt.Println()
			fmt.Println("tools:")
			fmt.Println(" 1. get server info")
			fmt.Println(" 2. clear webhooks")
			fmt.Println(" 3. fill channel list")
			fmt.Println(" 4. back")
			fmt.Println()

			switch question("select an option", []string{"1", "2", "3"}) {

			case "1":
				guild, err := dg.Guild(server.ID)
				if err != nil {

					fmt.Printf("[err]: unable to retrieve the guild object for the selected server...\n")
					fmt.Printf("       %v\n", err)
					os.Exit(1)

				}

				hooks, err := dg.GuildWebhooks(guild.ID)
				if err != nil {

					fmt.Printf("[err]: unable to retrieve webhooks for the selected server...\n")
					fmt.Printf("       %v\n", err)
					os.Exit(1)

				}

				channels, err = dg.GuildChannels(guild.ID)
				if err != nil {

					fmt.Printf("[err]: unable to retrieve the channels for the selected server...\n")
					fmt.Printf("       %v\n", err)
					os.Exit(1)

				}

				owner, err := dg.User(guild.OwnerID)
				if err != nil {

					fmt.Printf("[err]: unable to retrieve the guild owner for the selected server...\n")
					fmt.Printf("       %v\n", err)
					os.Exit(1)

				}

				fmt.Println()
				fmt.Printf("guild info for %s:\n", guild.Name)
				fmt.Printf("  channel count: %d\n", len(channels))
				fmt.Printf("  webhook count: %d\n", len(hooks))
				fmt.Printf("  member count: %d\n", guild.MemberCount)
				fmt.Printf("  owner:\n")
				fmt.Printf("    id: %s\n", owner.ID)
				fmt.Printf("    name: %s\n", owner.String())
				fmt.Printf("    bot: %v\n", owner.Bot)
				fmt.Printf("    verified: %v\n", owner.Verified)
				fmt.Printf("    avatar: %s\n", owner.AvatarURL(""))
				fmt.Println()

			case "2":
				hooks, err := dg.GuildWebhooks(server.ID)
				if err != nil {

					fmt.Printf("[err]: unable to retrieve webhooks for the selected server...\n")
					fmt.Printf("       %v\n", err)
					os.Exit(1)

				}

				for _, hook := range hooks {

					err := dg.WebhookDelete(hook.ID)
					if err != nil {

						fmt.Printf("[err]: unable to delete webhook...\n")
						fmt.Printf("       %v\n", err)

					}

				}

			case "3":
				channels, err = dg.GuildChannels(server.ID)
				if err != nil {

					fmt.Printf("[err]: unable to retrieve the channels for the selected server...\n")
					fmt.Printf("       %v\n", err)
					os.Exit(1)

				}

				if len(channels) == 100 {

					fmt.Printf("this server already has 100 channels\n")

				}

				for x := len(channels); x <= 100; x++ {

					_, err := dg.GuildChannelCreate(server.ID, "spam", "text")
					if err != nil {

						fmt.Printf("[err]: unable to create channel...\n")
						fmt.Printf("       %v\n", err)
						os.Exit(1)

					}

				}

				fmt.Printf("finished adding %d channel(s)...\n", 100-len(channels))

			}

		case "3":
			os.Exit(0)

		}
	}
}
