package main

import (
	"fmt"
	"github.com/gempir/go-twitch-irc/v3"
	"os"
	"path"
	"sync"
)

type Archiver struct {
	Channels  []string
	DirPrefix string

	mut       sync.RWMutex
	UserFiles map[string]*os.File
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: ncw <channelName> <anotherChannelName> ...\n")
		os.Exit(1)
	}

	a := Archiver{
		Channels:  os.Args[1:],
		DirPrefix: "./out/",
		mut:       sync.RWMutex{},
		UserFiles: make(map[string]*os.File),
	}
	err := a.setupDisk()
	if err != nil {
		fmt.Printf("failed to setup archiver: %s\n", err)
		return
	}

	client := twitch.NewAnonymousClient()
	client.OnPrivateMessage(a.writeMessageToDisk)
	client.Join(os.Args[1:]...)
	err = client.Connect()
	if err != nil {
		panic(err)
	}
}

func (a *Archiver) setupDisk() error {
	for _, c := range a.Channels {
		err := os.MkdirAll(path.Join(a.DirPrefix, c), 0775)
		if err != nil {
			return fmt.Errorf("failed to setup %s on disk: %w", c, err)
		}
	}
	return nil
}

func (a *Archiver) writeMessageToDisk(message twitch.PrivateMessage) {
	a.mut.RLock()
	f, exists := a.UserFiles[message.User.ID]
	a.mut.RUnlock()
	if !exists {
		a.mut.Lock()
		userFile, err := os.Create(path.Join(a.DirPrefix, message.Channel, message.User.Name))
		if err != nil {
			fmt.Printf("failed to setup file for %s under %s: %s\n", message.User.ID, path.Join(a.DirPrefix, message.Channel, message.User.Name), err)
			return
		}

		a.UserFiles[message.User.ID] = userFile
		a.mut.Unlock()
		f = userFile
	}

	_, err := fmt.Fprintf(f, "%s %s %s\n", message.Time, message.User.Name, message.Message)
	if err != nil {
		fmt.Printf("failed to write message to file: %s\n", err)
		return
	}
}
