package bcast

import (
	"fmt"
	"net"
	"time"
)

func MainKLoop(port int) {

	// 1. Start the Server in a goroutine
	go reciver(port)

	// Give the server a tiny bit of time to start listening
	time.Sleep(100 * time.Millisecond)

	// 2. Run the Client
	//sender(port)
}

// server listens for incoming UDP packets
func reciver(port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Println("Receiver faen error:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		fmt.Println("Receiver error:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Receiver started on %s\n", addr)

	buffer := make([]byte, 1024)
	for {
		// ReadFromUDP blocks until a packet arrives
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		fmt.Printf("[%s] Received: %s\n", remoteAddr, string(buffer[:n]))
	}
}

// client sends a packet and waits for a response
func sender(port int) {
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port))

	// DialUDP identifies the destination, but doesn't "connect"
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Sender error:", err)
		return
	}
	defer conn.Close()

	messages := []string{"Alpha", "Bravo", "Charlie", "Delta"}

	for _, m := range messages {
		fmt.Printf("Sender: Sending %s...\n", m)
		_, err := conn.Write([]byte(m))
		if err != nil {
			fmt.Println("Write error:", err)
		}
		time.Sleep(500 * time.Millisecond) // Space them out so it's readable
	}
}
