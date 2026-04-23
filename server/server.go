package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

// private
type stat struct {
	sync.Mutex
	count int
}

func (s *stat) Increment() int {
	s.Lock()
	defer s.Unlock()
	s.count++
	return s.count
}
func (s *stat) Decrement() int {
	s.Lock()
	defer s.Unlock()
	s.count--
	return s.count
}

// global
var activeClients stat

func readCommand(c net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}
func handleClient(conn net.Conn) {
	// deferred block: Even if something goes wrong, the counter will
	// be updated properly and connection will be closed
	defer func() {
		conn.Close()
		ct := activeClients.Decrement()
		log.Printf("Client disconnected. Total active: %d", ct)
	}()
	// This loop allows the client to send multiple commands in one session
	for {
		cmd, err := readCommand(conn)
		if err != nil {
			if err != io.EOF {
				// epoll notification: client hung up
				log.Println("Read error:", err)
			}
			break // Exit the inner loop to handle disconnect
		}

		log.Printf("Received: %q", cmd)
		if err = respond(cmd, conn); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func RunTCPServer(host string, port int) {
	addr := host + ":" + strconv.Itoa(port)
	log.Println("Starting Async TCP server on", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Failed to bind to port:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		newCount := activeClients.Increment()
		log.Printf("Client connected [%s]. Total active: %d", conn.RemoteAddr(), newCount)
		go handleClient(conn)
	}

}
