package netlib

import (
	"fmt"
	"net"
	"testing"
)

// go test -run "^TestServerConnection$" -v
func TestServerConnection(t *testing.T) {
	// listener, err := runTestServer()
	// if err != nil {
	// 	t.Fatalf("Failed to start server: %v", err)
	// }
	// defer listener.Close()
	listener, err := Sock()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	done := make(chan struct{})
	defer close(done)
	go Accept(listener, handleConnection, done)

	port := listener.Addr().(*net.TCPAddr).Port
	address := fmt.Sprintf("localhost:%d", port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send data
	message := "hello"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to write to connection: %v", err)
	}

	// Receive data
	buffer := make([]byte, len(message))
	_, err = conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read from connection: %v", err)
	}

	if string(buffer) != message {
		t.Errorf("Expected %q, got %q", message, string(buffer))
	}
}
