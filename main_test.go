package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockReader struct {
	mock.Mock
	text string
}

func (r *MockReader) Read(buffer []byte) (int, error) {
	bytes := []byte(r.text)
	n := len(bytes)
	copy(buffer, bytes)
	return n, nil
}

func TestRead(t *testing.T) {
	mockReader := &MockReader{}
	mockReader.text = "test"

	logReader := createLogReader()
	logReader.reader = mockReader

	text := logReader.read()
	assert.Equal(t, "test", text)

	mockReader.text = ""
	text = logReader.read()
	assert.Equal(t, "", text)
}

type MockHandler struct {
	mock.Mock
	handled bool
}

func (h *MockHandler) Handle(text string) {
	h.handled = true
}

func TestNotify(t *testing.T) {
	handler := &MockHandler{}

	logReader := createLogReader()
	logReader.expression = "Hello"
	logReader.handler = handler

	logReader.processLogs("test Hello")
	assert.True(t, handler.handled)
}
