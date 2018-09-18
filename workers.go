/*

workers.go -
file that contains goroutine workers

*/

package main

import (
	// internals
	"fmt"
	"time"
	// externals
	"github.com/bwmarrin/discordgo"
)

// the webhook worker
func webhookWorker(num int, name, message, avatarUrl string) {

	var (
		err       error
		tries     int
		sleepTime time.Duration
	)

	channel := channels[num/10]

	webhook, err := dg.WebhookCreate(channel.ID, name, "")
	if err != nil {

		// webhook recreation loop thing
		for tries = 0; tries <= 3; tries++ {

			webhook, err = dg.WebhookCreate(channel.ID, name, "")
			if err != nil {

				fmt.Printf("[err]: unable to recreate webhook...\n")
				fmt.Printf("       %v\n", err)

				if tries == 3 {

					return

				} else {

					sleepTime, err = time.ParseDuration(fmt.Sprintf("%ds", 30*tries))
					if err != nil {

						fmt.Printf("[err]: unable to parse duration...\n")
						fmt.Printf("       %v\n", err)
						return

					}

					time.Sleep(sleepTime)

				}
				continue

			}
			break

		}

	}

	for {

		err = dg.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
			Content:   message,
			AvatarURL: avatarUrl,
		})
		if err != nil {

			if err.(*discordgo.RESTError).Response.StatusCode == 404 {

				// webhook recreation loop thing
				for tries = 0; tries <= 3; tries++ {

					webhook, err = dg.WebhookCreate(channel.ID, name, "")
					if err != nil {

						fmt.Printf("[err]: unable to recreate webhook...\n")
						fmt.Printf("       %v\n", err)

						if tries == 3 {

							return

						} else {

							sleepTime, err = time.ParseDuration(fmt.Sprintf("%ds", 30*tries))
							if err != nil {

								fmt.Printf("[err]: unable to parse duration...\n")
								fmt.Printf("       %v\n", err)
								return

							}

							time.Sleep(sleepTime)

						}
						continue

					}
					break

				}

			} else if err.(*discordgo.RESTError).Response.StatusCode == 500 {

				fmt.Printf("[err]: an internal server error occured...\n")
				continue

			} else {

				fmt.Printf("[err]: unknown error from api...\n")
				fmt.Printf("       %v\n", err)
				return

			}

		}

	}

}
