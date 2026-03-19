package main

import (
	"Network-go/network/bcast"
	"flag"
	"fmt"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//
//	will be received as zero-values.
type HelloMsg struct {
	Message string
	Iter    int
}

func main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)

	//helloTx3 := make(chan HelloMsg)
	helloRx3 := make(chan HelloMsg)

	go bcast.Transmitter(16570, helloTx)
	go bcast.Receiver(16570, helloRx)

	// go bcast.Transmitter(16570, helloTx3)
	// go bcast.Receiver(16570, helloRx3)

	//go bcast.MainKLoop(16569 + 1000)

	time.Sleep(time.Second * 1)
	c := 1
	go func() {
		for {
			c += 1
			helloTx <- HelloMsg{Message: "Kuken", Iter: c}
			time.Sleep(time.Second)
		}
	}()

	//helloTx3 <- HelloMsg{Message: "Kuken 3", Iter: 1}

	for {
		select {
		case ans := <-helloRx:
			fmt.Println("Bober kurwa", ans)
		case ans2 := <-helloRx3:
			fmt.Println("Bober korwa faen: ", ans2)

		}

	}

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)

}

// if id == "" {
// 	localIP, err := localip.LocalIP()
// 	if err != nil {
// 		fmt.Println(err)
// 		localIP = "DISCONNECTED"
// 	}
// 	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
// }

// // We make a channel for receiving updates on the id's of the peers that are
// //  alive on the network
// peerUpdateCh := make(chan peers.PeerUpdate)
// // We can disable/enable the transmitter after it has been started.
// // This could be used to signal that we are somehow "unavailable".
// peerTxEnable := make(chan bool)
// go peers.Transmitter(15647, id, peerTxEnable)
// go peers.Receiver(15647, peerUpdateCh)

// // We make channels for sending and receiving our custom data types
// helloTx := make(chan HelloMsg)
// helloRx := make(chan HelloMsg)
// // ... and start the transmitter/receiver pair on some port
// // These functions can take any number of channels! It is also possible to
// //  start multiple transmitters/receivers on the same port.
// go bcast.Transmitter(16569, helloTx)
// go bcast.Receiver(16569, helloRx)

// // The example message. We just send one of these every second.
// go func() {
// 	helloMsg := HelloMsg{"Hello from " + id, 0}
// 	for {
// 		helloMsg.Iter++
// 		helloTx <- helloMsg
// 		time.Sleep(1 * time.Second)
// 	}
// }()

// fmt.Println("Started")
// for {
// 	select {
// 	case p := <-peerUpdateCh:
// 		fmt.Printf("Peer update:\n")
// 		fmt.Printf("  Peers:    %q\n", p.Peers)
// 		fmt.Printf("  New:      %q\n", p.New)
// 		fmt.Printf("  Lost:     %q\n", p.Lost)

// 	case  <-helloRx:
// 		// fmt.Printf("Received: %#v\n", a)
// 	}
// }
