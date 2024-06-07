package netlib

import (
	"fmt"
	"syscall"
)

const (
	maxConnections = 1024
	bufferSize     = 1024
)

type connection struct {
	fd        int
	readBuff  []byte
	writeBuff []byte
}

func Syscall_Sock() (syscall.Handle, int, error) {
	// syscall.AF_INET is the address family for IPv4 addresses
	// syscall.SOCK_STREAM is the socket type for a stream-oriented socket
	// 0 means that the default protocol for the specified address family and socket type will be used TCP
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("Failed to create socket:", err)
		return syscall.InvalidHandle, -1, err
	}
	//defer syscall.Close(fd)

	addr := &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}

	err = syscall.Bind(fd, addr)
	if err != nil {
		fmt.Println("Failed to bind socket:", err)
		return syscall.InvalidHandle, -1, err
	}

	err = syscall.Listen(fd, 10)
	if err != nil {
		fmt.Println("Failed to listen on socket:", err)
		return syscall.InvalidHandle, -1, err
	}
	fmt.Println("Listening on 127.0.0.1:8080")
	return fd, 0, nil
}

func Syscall_Accept(fd syscall.Handle) {
	for {
		connFd, _, err := syscall.Accept(fd)
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		fmt.Println("Accepted connection with file descriptor:", connFd)
		go syscall_handleConnection(connFd)
	}
}

func syscall_handleConnection(fd syscall.Handle) {
	defer syscall.Close(fd)

	buff := make([]byte, bufferSize)
	n, err := syscall.Read(fd, buff)
	if err != nil {
		fmt.Println("Failed to read from connection:", err)
		return
	}
	fmt.Println("Received data:", string(buff[:n]))

	_, err = syscall.Write(fd, []byte("Connetion is accepted ..."))
	if err != nil {
		fmt.Println("Failed to write to connection:", err)
		return
	}
}
