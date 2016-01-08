package main

import (
	"github.com/kardianos/osext"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
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
	path, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	handler := NewHandler()
	handler.Add(SoundNotification(path))
	handler.Add(UINotification(*expr))

	logProcessor.SetHandler(handler)
	logProcessor.start()
}
