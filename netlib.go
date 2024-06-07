package netlib

import (
	"fmt"
	"io"
	"net"
	"os"
)

func Sock() (net.Listener, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create listener: %v\n", err)
		return nil, err
	}
	fmt.Println("Listening on", listener.Addr().String())
	return listener, nil
}

type handlerConnFunc func(c net.Conn)

func Accept(listener net.Listener, cb handlerConnFunc, done chan struct{}) {
	//defer listener.Close()
	for {
		select {
		case <-done:
			fmt.Println("Shutdown signal received, stopping server...")
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				// listener was closed, exit the loop
				if err == net.ErrClosed {
					return
				}
				fmt.Fprintf(os.Stderr, "Failed to accept connection: %v\n", err)
				continue
			}
			fmt.Println("Accepted connection from", conn.RemoteAddr().String())
			go cb(conn)
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Failed to read from connection:", err)
		return
	}

	fmt.Println("Received data:", string(buffer[:n]))

	_, err = conn.Write(buffer[:n])
	if err != nil {
		fmt.Println("Failed to write to connection:", err)
		return
	}
}

func handleEcho(c net.Conn) {
	defer c.Close()
	io.Copy(c, c)
}

// func runTestServer(port string) (*net.Listener, error) {
// 	listener, err := net.Listen("tcp", "localhost:"+port)
// 	if err != nil {
// 		return nil, err
// 	}
// 	go func() {
// 		for {
// 			conn, err := listener.Accept()
// 			if err != nil {
// 				return
// 			}
// 			go func(c net.Conn) {
// 				defer c.Close()
// 				io.Copy(c, c)
// 			}(conn)
// 		}
// 	}()
// 	return &listener, nil
// }

func AcceptOne(listener net.Listener, cb handlerConnFunc, done chan struct{}) {
	conn, err := listener.Accept()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to accept connection: %v\n", err)
		close(done)
		return
	}
	go cb(conn)
	close(done)
}

func runTestServer() (net.Listener, error) {
	listener, err := Sock()
	if err != nil {
		return nil, err
	}

	done := make(chan struct{})
	go func() {
		AcceptOne(listener, handleConnection, done)
	}()

	return listener, nil
}
