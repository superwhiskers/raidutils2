/*

workers.go -
file that contains goroutine workers

*/

package main

import (
	// internals
	"fmt"
	// externals
	"github.com/bwmarrin/discordgo"
)

// the webhook worker
func webhookWorker(num int, respawn, cycle bool, name, message, avatarUrl string) {

	var (
		err     error
		webhook *discordgo.Webhook
	)

	channel := channels[num/10]

	if cycle != true {

		webhook, err = createWebhook(channel, name)
		if err != nil {

			fmt.Printf("[err]: unable to create webhook...\n")
			fmt.Printf("       %v\n", err)
			return

		}

	}

	for {

		if cycle == true {

			webhook, err = createWebhook(channel, name)
			if err != nil {

				fmt.Printf("[err]: unable to create webhook...\n")
				fmt.Printf("       %v\n", err)
				return

			}

		}

		err = dg.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
			Content:   message,
			AvatarURL: avatarUrl,
		})

		if err != nil {

			switch t := err.(type) {

			case *discordgo.RESTError:
				if err.(*discordgo.RESTError).Response.StatusCode == 404 {

					if cycle == true {

						fmt.Printf("[err]: webhook was deleted...\n")
						continue

					}

					webhook, err = createWebhook(channel, name)
					if err != nil {

						fmt.Printf("[err]: unable to create webhook...\n")
						fmt.Printf("       %v\n", err)
						return

					}

				} else if err.(*discordgo.RESTError).Response.StatusCode == 500 {

					fmt.Printf("[err]: an internal server error occurred...\n")
					continue

				} else if err.(*discordgo.RESTError).Response.StatusCode == 400 {

					fmt.Printf("[err]: a bad request error occurred...\n")
					continue

				} else {

					fmt.Printf("[err]: unknown error from api...\n")
					fmt.Printf("       %v\n", err)
					return

				}

			default:
				fmt.Printf("[err]: unknown error type from webhook execution %T\n", t)

			}

		}

		if cycle == true {

			// there's really no point in getting the error from this but i'll log it anyways
			err = dg.WebhookDelete(webhook.ID)
			if err != nil {

				fmt.Printf("[err]: unable to delete webhook...")
				fmt.Printf("       %v\n", err)

			}

		}

	}

}
