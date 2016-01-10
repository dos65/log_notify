package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	filePath = kingpin.Flag("file", "Path to file").Short('f').String()
	expr     = kingpin.Flag("expr", "Expression").Short('e').Required().String()
)

func main() {
	kingpin.Parse()

	var logProcessor *LogProcessor
	if *filePath == "" {
		logProcessor = NewStdinLogProcessor(*expr)
	} else {
		logProcessor = NewFileLogProcessor(*filePath, *expr)
	}

	handler := NewHandler()
	handler.Add(SoundNotification())
	handler.Add(UINotification())

	logProcessor.SetHandler(handler)
	logProcessor.start()
}
