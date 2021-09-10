/*

raid-utils.go -
the main file of the raid-utils rewrite

*/

package main

import (
	// internals
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	// externals
	"github.com/bwmarrin/discordgo"
)

// variables that are used to perform the raids
var (
	server   *discordgo.UserGuild
	channels []*discordgo.Channel
	dg       *discordgo.Session
	err      error

	searchMut = &sync.Mutex{}
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

		fmt.Printf("[err]: %s is not a number...\n", serverIndStr)
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	if len(servers) > serverIndInt && serverIndInt > -1 {

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
		fmt.Println(" 2. hybrid/distributed mode (experimental)")
		fmt.Println(" 3. open tool menu")
		fmt.Println(" 4. change server")
		fmt.Println(" 5. exit")
		fmt.Println()

		switch question("select an option", []string{"1", "2", "3", "4", "5"}) {

		case "1":
			menuOptRaid()

		case "2":
			initHybridModeServer()

		case "3":
			fmt.Println()
			fmt.Println("tools:")
			fmt.Println(" 1. get server info")
			fmt.Println(" 2. clear webhooks")
			fmt.Println(" 3. fill channel list")
			fmt.Println(" 4. retrieve invite code")
			fmt.Println(" 5. leave every server")
			fmt.Println(" 6. leave the current server")
			fmt.Println(" 7. go back")
			fmt.Println()

			switch question("select an option", []string{"1", "2", "3", "4", "5", "6", "7"}) {

			case "1":
				menuOptServerInfo()

			case "2":
				menuOptDeleteWebhooks()

			case "3":
				menuOptAddChannels()

			case "4":
				menuOptGetInvite()

			case "5":
				menuOptLeaveAll()

			case "6":
				menuOptLeaveCurrent()

			}

		case "4":
			menuOptChangeServer()

		case "5":
			os.Exit(0)

		}
	}
}
