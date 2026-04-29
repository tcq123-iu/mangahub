// cmd/tcp-server/server.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

// Simple Message Protocol
type Hub struct {
	clients map[net.Conn]bool
	mutex   sync.Mutex
}

var hub = Hub{
	clients: make(map[net.Conn]bool),
}

func main() {
	port := "8081"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("MangaHub TCP Sync Server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		// Requirement: Connection handling with goroutines
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	hub.mutex.Lock()
	hub.clients[conn] = true
	hub.mutex.Unlock()

	log.Printf("New device synced: %s", conn.RemoteAddr())

	defer func() {
		hub.mutex.Lock()
		delete(hub.clients, conn)
		hub.mutex.Unlock()
		conn.Close()
		log.Printf("Device disconnected: %s", conn.RemoteAddr())
	}()

	// Read incoming broadcast messages from the API Server
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		broadcast(msg, conn)
	}
}

func broadcast(message string, sender net.Conn) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	for client := range hub.clients {
		// Send to everyone EXCEPT the sender (usually the API server)
		if client != sender {
			fmt.Fprintln(client, message)
		}
	}
}