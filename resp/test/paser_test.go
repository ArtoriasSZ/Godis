package test

import (
	"Godis/resp/parser"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestReadFull(t *testing.T) {
	msg := make([]byte, 0, 5)

	_, err := io.ReadFull(os.Stdin, msg)
	fmt.Println(err)
	fmt.Println(len(msg), cap(msg), msg)
	_, err = io.ReadFull(os.Stdin, msg)
	fmt.Println(err)
	fmt.Println(len(msg), cap(msg), msg)
}

// MockReader implements io.Reader for testing purposes.
type MockReader struct {
	data []byte
}

func (r *MockReader) Read(p []byte) (n int, err error) {
	copy(p, r.data)
	return len(r.data), io.EOF
}

// TestParser0 tests the parser0 function.
func TestParser0(t *testing.T) {
	// Mock data for testing
	mockData := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"

	// Create a mock reader with mock data
	mockReader := &MockReader{data: []byte(mockData)}

	// Create a channel for receiving payloads

	// Run the parser0 function in a goroutine
	ch := parser.ParseStream(mockReader)

	// Wait for the payload from the channel

	receivedPayload := <-ch
	fmt.Println(receivedPayload.Date)
	fmt.Println(receivedPayload.Err)
	time.Sleep(100 * time.Second)
}
