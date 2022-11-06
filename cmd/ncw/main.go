package main

import (
	"fmt"
	"github.com/gempir/go-twitch-irc/v3"
	"os"
	"os/signal"
)

func main() {
	if len(os.Args) == 0 {
		fmt.Printf("Usage: ncw <channelName> <anotherChannelName> ...")
		os.Exit(1)
	}

	var chatUserDB = make(map[string]struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("%+v\n", chatUserDB)
		os.Exit(1)
	}()

	client := twitch.NewAnonymousClient()
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Printf("%s: %s %s - %s\n", message.Channel, message.Time, message.User.Name, message.Message)
		_, exists := chatUserDB[message.User.ID]
		if exists {
			return
		}
		chatUserDB[message.User.ID] = struct{}{}
	})
	client.Join(os.Args[1:]...)
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
