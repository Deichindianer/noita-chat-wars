package main

import (
	"fmt"
	"github.com/gempir/go-twitch-irc/v3"
	"os"
	"os/signal"
)

func main() {
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
		fmt.Printf("%s %s - %s\n", message.Time, message.User.Name, message.Message)
		_, exists := chatUserDB[message.User.ID]
		if exists {
			return
		}
		chatUserDB[message.User.ID] = struct{}{}
	})
	client.Join("Zizaran")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
