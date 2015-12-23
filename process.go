package main

import (
	"fmt"
	"github.com/mgutz/ansi"
	"golang.org/x/exp/inotify"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
    sleepTime = 1000
	bytesRead = 512
)

type LogProcessor struct {
	path string
	expression string
	reader io.Reader
	inotify bool
    handler Handler
	buffer []byte
}

func createLogProcessor() *LogProcessor {
	buffer := make([]byte, bytesRead)
	return &LogProcessor{buffer: buffer}
}

func NewFileLogProcessor (path string, expression string) *LogProcessor {
    absPath := absPath(path)
	logProcessor := createLogProcessor()
	logProcessor.path = absPath
	logProcessor.expression = expression
	logProcessor.inotify = true
	logProcessor.reader = createReader(absPath)
	return logProcessor
}

func NewStdinLogProcessor ( expression string) *LogProcessor {
	logProcessor := createLogProcessor()
	logProcessor.expression = expression
	logProcessor.reader = os.Stdin
	return logProcessor
}

func (l * LogProcessor) SetHandler(handler Handler) {
	l.handler = handler
}

func absPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	return abs
}

func createReader(path string) io.Reader {
    file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    file.Seek(0, 2)
    return file
}

func (l *LogProcessor) start() {
    if l.inotify {
        l.startInotify()
    } else {
        l.startDefault()
    }
}

func (l *LogProcessor) startInotify() {
    watcher, err := inotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    err = watcher.Watch(l.path)
    if err != nil {
        log.Fatal(err)
    }
    for {
        select {
        case event := <-watcher.Event:
            if event.Mask == inotify.IN_MODIFY {
                l.read()
            }
        case err := <-watcher.Error:
            log.Fatal(err)
        }
    }
}

func (l *LogProcessor) startDefault() {
    for {
        l.read()
        time.Sleep(1000)
    }

}

func (l *LogProcessor) read() string {
    for {
        n, err := l.reader.Read(l.buffer)
        if err != nil {
            if err == io.EOF {
                return ""
            }
            log.Fatal(err)
        }

        if n > 0 {
			text := string(l.buffer[:n])
			if l.expression != "" {
				l.processLogs(text)
			}
			return text
        } else {
            return ""
        }
    }
}

func (l *LogProcessor) processLogs(text string) {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		match, _ := regexp.MatchString(l.expression, line)
		if match {
			l.handler.Handle(lines[i])
			lines[i] = ansi.Color(line, "red")
		}
	}
	fmt.Print(strings.Join(lines, "\n"))
}
