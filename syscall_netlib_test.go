package netlib

import (
	"os"
	"syscall"
	"testing"
)

var (
	FD syscall.Handle
)

func TestMain(m *testing.M) {
	exitCode := m.Run()

	fd, n, err := Syscall_Sock()
	if err != nil && n != 0 {
		panic(err)
	}
	FD = fd

	os.Exit(exitCode)
}

func Test_Plain(t *testing.T) {

	// fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	// if err != nil {
	// 	fmt.Println("Failed to connect to socket:", err)
	// 	return
	// }
	// fmt.Println("Socket created with file descriptor:", fd)

	t.Run("Testing ...", func(t *testing.T) {
		t.Log("started socket on addr:port http://localhost:8080")
		Syscall_Accept(FD)
	})
}
