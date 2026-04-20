package server

import (
	"io"
	"log"
	"net"
	"strconv"
)

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

func RunSyncTCPServer(host string, port int) {
	addr := host + ":" + strconv.Itoa(port)
	log.Println("Starting synchronous TCP server on", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Failed to bind to port:", err)
	}
	defer listener.Close()

	clients := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}

		clients++
		log.Printf("Client connected [%s]. Total active: %d", conn.RemoteAddr(), clients)

		// This loop allows the client to send multiple commands in one session
		for {
			cmd, err := readCommand(conn)
			if err != nil {
				if err != io.EOF {
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

		conn.Close()
		clients--
		log.Printf("Client disconnected. Total active: %d", clients)
	}
}
