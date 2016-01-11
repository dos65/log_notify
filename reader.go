package main

import (
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

type LogReader struct {
	path       string
	expression string
	reader     io.Reader
	inotify    bool
	handler    Handler
	buffer     []byte
	out        chan string
}

func createLogReader() *LogReader {
	buffer := make([]byte, bytesRead)
	out := make(chan string)
	return &LogReader{buffer: buffer, out: out}
}

func NewFileLogReader(path string, expression string) *LogReader {
	absPath := absPath(path)
	logReader := createLogReader()
	logReader.path = absPath
	logReader.expression = expression
	logReader.inotify = true
	logReader.reader = createReader(absPath)
	return logReader
}

func NewStdinLogReader(expression string) *LogReader {
	logProcessor := createLogReader()
	logProcessor.expression = expression
	logProcessor.reader = os.Stdin
	return logProcessor
}

func (l *LogReader) SetHandler(handler Handler) {
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

func (l *LogReader) start() {
	if l.inotify {
		l.startInotify()
	} else {
		l.startDefault()
	}
}

func (l *LogReader) startInotify() {
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
				l.tryRead()
			}
		case err := <-watcher.Error:
			log.Fatal(err)
		}
	}
}

func (l *LogReader) startDefault() {
	for {
		l.tryRead()
		time.Sleep(1000)
	}
}

func (l *LogReader) tryRead() {
	text := l.read()
	if text == "" {
		return
	}
	text = l.processLogs(text)
	l.out <- text
}

func (l *LogReader) read() string {
	for {
		n, err := l.reader.Read(l.buffer)
		if err != nil {
			if err == io.EOF {
				return ""
			}
			log.Fatal(err)
		}

		if n > 0 {
			return string(l.buffer[:n])
		} else {
			return ""
		}
	}
}

func (l *LogReader) processLogs(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		match, _ := regexp.MatchString(l.expression, line)
		if match {
			l.handler.Handle(lines[i])
			lines[i] = ansi.Color(line, "red")
		}
	}
	return strings.Join(lines, "\n")
}
