package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type botConfig struct {
	// Port is server port to be listened.
	Port string `envconfig:"PORT" default:"3000"`

	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// VerificationToken is used to validate interactive messages from slack.
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID" required:"true"`
}

func main() {
	var env botConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalln("[ERROR] Failed to process env var: %s", err)
	}

	api := slack.New(env.BotToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				handleMessageEvent(MessageContext{
					api, env, rtm,
				}, ev)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")

			default:
				// Take no action
			}
		}
	}
}
