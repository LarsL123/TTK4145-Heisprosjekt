package reelection_test

// For testing package reelection

import (
	"fmt"
	"net"
	"reelection"
	"testing"
)

func sendheartbeat() {

	addr, err := net.ResolveUDPAddr("udp", "localhost:31000")

	if err != nil {
		fmt.Print("Fant ingen adresse:", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		log.Fatal("Rompe ass vi fikk error ", err)
	}
	defer conn.Close()
	//buffer := make([]byte, 1024)
	//response := fmt.Print("Rompami")
	for {
		_, err = conn.Write([]byte("rompami"))
		// n, serverAddr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			log.Printf("Fittefaenhælvete: %s", err)
			continue
		}

		//message := string(buffer[:n])

		time.Sleep(10000 * time.Millisecond)
	}
}


func receiveheartbeat() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:31000")

	if err != nil {
		fmt.Print("Fant ingen adresse:", err)
	}

	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		log.Fatal("Rompe ass vi fikk error ", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, serverAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Fittefaen")
			continue
		}

		message := string(buffer[:n])

		fmt.Printf("Vi fikk melding fra %s: %s\n", serverAddr, message)
	}

}


///////////////////////////// 