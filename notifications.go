package main

import (
	"os/exec"
	"path/filepath"
	"log"
)

type Notification func(string)

func UINotification(text string) Notification {
	return func(text string) {
		err := exec.Command("notify-send", "lognotify" , "\"" + text + "\" handled").Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func SoundNotification(path string) Notification {
	return func(text string) {
		exec.Command("ogg123", filepath.Join(path, "ring.ogg")).Run()
	}
}

