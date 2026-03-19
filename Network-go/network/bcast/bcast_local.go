package bcast

import (
	"Network-go/network/conn"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

// Encodes received values from `chans` into type-tagged JSON, then broadcasts
// it on `port`
func TransmitterLocal(port int, chans ...interface{}) {
	checkArgs(chans...)
	typeNames := make([]string, len(chans))
	selectCases := make([]reflect.SelectCase, len(typeNames))
	for i, ch := range chans {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
		typeNames[i] = reflect.TypeOf(ch).Elem().String()
	}

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	//Local
	addrLocal, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port+1000))
	connLocal, _ := net.DialUDP("udp", nil, addrLocal)

	for {
		chosen, value, _ := reflect.Select(selectCases)
		jsonstr, _ := json.Marshal(value.Interface())
		ttj, _ := json.Marshal(typeTaggedJSON{
			TypeId: typeNames[chosen],
			JSON:   jsonstr,
		})
		if len(ttj) > bufSize {
			panic(fmt.Sprintf(
				"Tried to send a message longer than the buffer size (length: %d, buffer size: %d)\n\t'%s'\n"+
					"Either send smaller packets, or go to network/bcast/bcast.go and increase the buffer size",
				len(ttj), bufSize, string(ttj)))
		}
		conn.WriteTo(ttj, addr)
		connLocal.Write(ttj)

	}
}

// Matches type-tagged JSON received on `port` to element types of `chans`, then
// sends the decoded value on the corresponding channel
func ReceiverLocal(port int, chans ...interface{}) {
	checkArgs(chans...)
	chansMap := make(map[string]interface{})
	for _, ch := range chans {
		chansMap[reflect.TypeOf(ch).Elem().String()] = ch
	}

	//Recive Local sjef
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port+1000))
	if err != nil {
		fmt.Println("Receiver faen error:", err)
		return
	}
	connLocal, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Receiver error:", err)
		return
	}
	defer connLocal.Close()

	go func() {
		buffer := make([]byte, 1024)
		for {
			// ReadFromUDP blocks until a packet arrives
			n, remoteAddr, err := connLocal.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Error reading:", err)
				continue
			}
			fmt.Println(remoteAddr)
			var ttj typeTaggedJSON
			json.Unmarshal(buffer[0:n], &ttj)
			ch, ok := chansMap[ttj.TypeId]
			if !ok {
				continue
			}
			v := reflect.New(reflect.TypeOf(ch).Elem())
			json.Unmarshal(ttj.JSON, v.Interface())
			reflect.Select([]reflect.SelectCase{{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(ch),
				Send: reflect.Indirect(v),
			}})

		}
	}()

	var buf [bufSize]byte
	conn := conn.DialBroadcastUDP(port)
	for {
		n, _, e := conn.ReadFrom(buf[0:])
		if e != nil {
			fmt.Printf("bcast.Receiver(%d, ...):ReadFrom() failed: \"%+v\"\n", port, e)
		}

		var ttj typeTaggedJSON
		json.Unmarshal(buf[0:n], &ttj)
		ch, ok := chansMap[ttj.TypeId]
		if !ok {
			continue
		}
		v := reflect.New(reflect.TypeOf(ch).Elem())
		json.Unmarshal(ttj.JSON, v.Interface())
		reflect.Select([]reflect.SelectCase{{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: reflect.Indirect(v),
		}})
	}
}
