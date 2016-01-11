package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	filePath = kingpin.Flag("file", "Path to file").Short('f').String()
	expr     = kingpin.Flag("expr", "Expression").Short('e').Required().String()
)

func main() {
	kingpin.Parse()

	var logProcessor *LogReader
	if *filePath == "" {
		logProcessor = NewStdinLogReader(*expr)
	} else {
		logProcessor = NewFileLogReader(*filePath, *expr)
	}

	handler := NewHandler()
	handler.Add(SoundNotification())
	handler.Add(UINotification())

	logProcessor.SetHandler(handler)
	go logProcessor.start()
	for output := range logProcessor.out {
		fmt.Print(output)
	}
}
