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

func (r *MockReader) Read(buffer []byte) (int , error) {
	bytes := []byte(r.text)
	n := len(bytes)
	copy(buffer, bytes)
	return n, nil
}

func TestRead(t *testing.T) {
	mockReader := &MockReader{}
	mockReader.text = "test"

	logProcessor := createLogProcessor()
	logProcessor.reader = mockReader

	text := logProcessor.read()
	assert.Equal(t, "test", text)

	mockReader.text = ""
	text = logProcessor.read()
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
	mockReader := &MockReader{}
	mockReader.text = "test Hello"

	handler := &MockHandler{}

	logProcessor := createLogProcessor()
	logProcessor.reader = mockReader
	logProcessor.expression = "Hello"
	logProcessor.handler = handler

	text := logProcessor.read()
	t.Log("text:" + text)
	assert.True(t, handler.handled)
}
