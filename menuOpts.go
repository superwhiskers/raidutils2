/*

menuOpts.go -
functions that run menu options

*/

package main

import (
	// internals
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	// externals
	"github.com/bwmarrin/discordgo"
)

// raid menu option code
func menuOptRaid() {

	code := strconv.Itoa(randint(100000, 999999))
	fmt.Printf("warning~~!\n")
	fmt.Printf("you could potentially get banned from discord or\n")
	fmt.Printf("the server you are using this in with this feature...\n")
	fmt.Printf("enter the following code to perform this action: %s\n", code)

	if code == input(": ") {

		name := question("what should the webhooks be named?", []string{})
		message := question("what should they spam?", []string{})
		avatarUrl := question("what should the avatar url be?", []string{})
		cycleWebhooks := convertYNToBool[question("should a new webhook be used for every new spam message?", []string{"yes", "no"})]
		limitWebhooks := question("should there be a webhook cap?", []string{"yes", "no"})

		var hooks int
		if limitWebhooks == "no" {

			fmt.Println("calculating webhook count...")

		} else {

			hooksTmp := question("what should they be capped at?", []string{})
			hooks, err = strconv.Atoi(hooksTmp)
			if err != nil {

				fmt.Printf("%s is not a number...\n", hooksTmp)
				return

			}

		}

		allChannels, err := dg.GuildChannels(server.ID)
		if err != nil {

			fmt.Printf("[err]: unable to retrieve channels for the selected server... (returning to menu)\n")
			fmt.Printf("       %v\n", err)
			return

		}

		channels = []*discordgo.Channel{}

		for _, channel := range allChannels {

			if channel.Type == discordgo.ChannelTypeGuildText {

				channels = append(channels, channel)

			}

		}

		if limitWebhooks == "no" {

			hooks = len(channels) * 10

		}

		fmt.Printf("approximately %d webhooks will be spawned\n", hooks)
		if question("are you sure you want to proceed?", []string{"yes", "no"}) == "no" {

			return

		}

		fmt.Printf("spawned 0/%d workers...", hooks)

		for w := 1; w <= hooks; w++ {

			go webhookWorker(w-1, cycleWebhooks, name, message, avatarUrl)

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

}

// server info menu option
func menuOptServerInfo() {

	guild, err := dg.Guild(server.ID)
	if err != nil {

		fmt.Printf("[err]: unable to retrieve the guild object for the selected server... (returning to menu)\n")
		fmt.Printf("       %v\n", err)
		return

	}

	hooks, err := dg.GuildWebhooks(guild.ID)
	if err != nil {

		fmt.Printf("[err]: unable to retrieve webhooks for the selected server... (continuing anyways)\n")
		fmt.Printf("       %v\n", err)

	}

	channels, err = dg.GuildChannels(guild.ID)
	if err != nil {

		fmt.Printf("[err]: unable to retrieve the channels for the selected server... (continuing anyways)\n")
		fmt.Printf("       %v\n", err)

	}

	owner, err := dg.User(guild.OwnerID)
	if err != nil {

		fmt.Printf("[err]: unable to retrieve the guild owner for the selected server... (continuing anyways)\n")
		fmt.Printf("       %v\n", err)

	}

	invites, err := dg.GuildInvites(guild.ID)
	if err != nil {

		fmt.Printf("[err]: unable to retrieve the invites for the selected server... (continuing anyways)\n")
		fmt.Printf("       %v\n", err)

	}

	var invite string
	if len(invites) != 0 {

		invite = invites[0].Code

	} else {

		invite = "none"

	}

	fmt.Println()
	fmt.Printf("guild info for \"%s\":\n", guild.Name)
	fmt.Printf("  region: %s\n", guild.Region)
	fmt.Printf("  channel count: %d\n", len(channels))
	fmt.Printf("  webhook count: %d\n", len(hooks))
	fmt.Printf("  member count: %d\n", guild.MemberCount)
	fmt.Printf("  first available invite: %s\n", invite)
	fmt.Printf("  owner:\n")
	fmt.Printf("    id: %s\n", owner.ID)
	fmt.Printf("    name: %s\n", owner.String())
	fmt.Printf("    bot: %v\n", owner.Bot)
	fmt.Printf("    verified: %v\n", owner.Verified)
	fmt.Printf("    avatar: %s\n", owner.AvatarURL(""))
	fmt.Println()

}

// menu option for webhook deletion
func menuOptDeleteWebhooks() {

	hooks, err := dg.GuildWebhooks(server.ID)
	if err != nil {

		fmt.Printf("[err]: unable to retrieve webhooks for the selected server... (returning to menu)\n")
		fmt.Printf("       %v\n", err)
		return

	}

	fmt.Printf("deleting webhook 0/%d...", len(hooks))

	for i, hook := range hooks {

		fmt.Printf("\rdeleting webhook %d/%d...", i+1, len(hooks))

		err := dg.WebhookDelete(hook.ID)
		if err != nil {

			fmt.Printf("\n[err]: unable to delete webhook...\n")
			fmt.Printf("       %v\n", err)

		}

	}

	fmt.Printf("\nsuccessfully deleted %d webhooks\n", len(hooks))

}

// menu option that fills the channel list
func menuOptAddChannels() {

	var channelCount int

	channels, err = dg.GuildChannels(server.ID)
	if err != nil {

		fmt.Printf("[err]: unable to retrieve the channels for the selected server... (returning to menu)\n")
		fmt.Printf("       %v\n", err)
		return

	}

	fmt.Println()
	fmt.Println("select a channel count to reach:")
	fmt.Println("  1. 500 (max)")
	fmt.Println("  2. 400")
	fmt.Println("  3. 300")
	fmt.Println("  4. 200")
	fmt.Println("  5. 100")
	fmt.Println()

	switch question("how many channels would you like the server to have?", []string{"1", "2", "3", "4", "5"}) {

	case "1":
		channelCount = 500

	case "2":
		channelCount = 400

	case "3":
		channelCount = 300

	case "4":
		channelCount = 200

	case "5":
		channelCount = 100

	}

	if len(channels) == channelCount {

		fmt.Printf("this server already has %d channels!\n", channelCount)

	}

	fmt.Printf("adding channel 0/%d...", channelCount-len(channels))

	for x := len(channels); x < channelCount; x++ {

		fmt.Printf("\radding channel %d/%d...", (x+1)-len(channels), channelCount-len(channels))

		_, err := dg.GuildChannelCreate(server.ID, strings.Join([]string{randalphastring(randint(5, 12)), strconv.Itoa(randint(10000, 999999999))}, ""), discordgo.ChannelTypeGuildText)
		if err != nil {

			fmt.Printf("\n[err]: unable to create channel... (returning to menu)\n")
			fmt.Printf("       %v\n", err)
			return

		}

	}

	fmt.Printf("\nfinished adding %d channel(s)...\n", channelCount-len(channels))

}

// menu option that gets the server invite
func menuOptGetInvite() {

	invites, err := dg.GuildInvites(server.ID)
	if err != nil {

		fmt.Printf("[err]: unable to retrieve the invites for the selected server... (continuing anyways)\n")
		fmt.Printf("       %v\n", err)

	}

	var invite string
	if len(invites) == 0 {

		allChannels, err := dg.GuildChannels(server.ID)
		if err != nil {

			fmt.Printf("[err]: unable to retrieve the channels for the selected server... (returning to menu)\n")
			fmt.Printf("       %v\n", err)
			return

		}

		var channel *discordgo.Channel

		for _, tmpChannel := range allChannels {

			if tmpChannel.Type == discordgo.ChannelTypeGuildText {

				channel = tmpChannel
				break

			}

		}

		inviteTmp, err := dg.ChannelInviteCreate(channel.ID, discordgo.Invite{
			MaxAge:    0,
			MaxUses:   0,
			Temporary: false,
		})
		if err != nil {

			fmt.Printf("[err]: unable to create an invite for the selected server... (returning to menu)\n")
			fmt.Printf("       %v\n", err)
			return

		}

		invite = inviteTmp.Code

	} else {

		invite = invites[0].Code

	}

	fmt.Println()
	fmt.Printf("invite code: %s", invite)
	fmt.Println()

}

// menu option that changes the target server
func menuOptChangeServer() {

	servers, err := dg.UserGuilds(100, "", "")
	if err != nil {

		fmt.Printf("[err]: could not retrieve the guilds for the bot... (returning to menu)\n")
		fmt.Printf("       %v\n", err)
		return

	}

	for i, s := range servers {

		fmt.Printf("%d: %s\n", i, s.Name)

	}

	serverIndStr := question("select a server", []string{})

	serverIndInt, err := strconv.Atoi(serverIndStr)
	if err != nil {

		fmt.Printf("[err]: %s is not a number... (returning to menu)\n", serverIndStr)
		fmt.Printf("       %v\n", err)
		return

	}

	if len(servers) > serverIndInt && serverIndInt > -1 {

		server = servers[serverIndInt]

	} else {

		fmt.Printf("[err]: %s is not in the server list... (returning to the menu)\n", serverIndStr)

	}

}

// makes the target user leave the currently selected server
func menuOptLeaveCurrent() {

	err = dg.GuildLeave(server.ID)
	if err != nil {

		fmt.Printf("\n[err]: unable to leave %s... (continuing anyways)\n", server.Name)
		fmt.Printf("       %v\n", err)

	}

	menuOptChangeServer()

}

// makes the target user leave every server that they're in
func menuOptLeaveAll() {

	servers, err := dg.UserGuilds(100, "", "")
	if err != nil {

		fmt.Printf("[err]: could not retrieve the guilds for the bot... (returning to menu)\n")
		fmt.Printf("       %v\n", err)
		return

	}

	length := len(servers)
	fmt.Printf("leaving server 0/%d...", length)

	for i, s := range servers {

		fmt.Printf("\rleaving server %d/%d...", i+1, length)

		err = dg.GuildLeave(s.ID)
		if err != nil {

			fmt.Printf("\n[err]: unable to leave %s... (continuing anyways)\n", s.Name)
			fmt.Printf("       %v\n", err)

		}

	}

	fmt.Printf("\n")

}
