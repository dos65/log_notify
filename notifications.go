package main

import (
	"github.com/kardianos/osext"
	"log"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type Notification interface {
	Notify(string)
}

type CommandNotification struct {
	name        string
	lockTimeout time.Duration
	getArgs     func(string) []string
	sync.Mutex
}

func (cn *CommandNotification) Notify(text string) {
	cn.Lock()
	err := exec.Command(cn.name, cn.getArgs(text)...).Run()
	if err != nil {
		log.Fatal(err)
	}
	if cn.lockTimeout > 0 {
		time.Sleep(cn.lockTimeout)
	}
	cn.Unlock()
}

func UINotification() Notification {
	notifySend := "notify-send"
	checkCmd(notifySend)
	getArgs := func(text string) []string {
		return []string{"lognotify", "\"" + text + "\" handled"}
	}
	return &CommandNotification{
		name:        notifySend,
		lockTimeout: 9000,
		getArgs:     getArgs}

}

func SoundNotification() Notification {
	ogg := "ogg123"
	checkCmd(ogg)
	path, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	oggPath := filepath.Join(path, "ring.ogg")
	return &CommandNotification{
		name: ogg,
		getArgs: func(text string) []string {
			return []string{oggPath}
		}}
}

func checkCmd(name string) {
	err := exec.Command("which", name).Run()
	if err != nil {
		log.Fatalf("Can not find %s.\nSee https://github.com/dos65/lognotify#install", name)
	}
}
