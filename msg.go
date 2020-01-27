package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

type MessageContext struct {
	client *slack.Client
	bconf  botConfig
	rtm    *slack.RTM
}

func handleMessageEvent(mctx MessageContext, ev *slack.MessageEvent) error {
	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", mctx.bconf.BotID)) {
		return nil
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]

	if len(m) == 0 {
		return fmt.Errorf("no user group provided")
	}

	members, err := mctx.client.GetUserGroupMembers(m[0])

	if err != nil {
		return err
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	index := r.Intn(len(members))
	msg := mctx.rtm.NewOutgoingMessage("selected "+members[index], ev.Channel)
	mctx.rtm.SendMessage(msg)

	return nil
}
