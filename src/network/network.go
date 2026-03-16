package network

import (
	"context"
	"elevatorproject/src/config"
	"fmt"
	"time"
)

// func (r *GenericSender[A, B]) SendAsyncWithAck(msg A) {
// 	r.mu.Lock()
// 	// cancel previous send if exists
// 	if r.cancelLast != nil {
// 		r.cancelLast()
// 	}
// 	ctx, cancel := context.WithCancel(context.Background())
// 	r.cancelLast = cancel
// 	r.lastUpdateNr = msg.GetUpdateNr()
// 	r.mu.Unlock()

// 	go func() {
// 		retryTicker := time.NewTicker(config.Cfg.AckRetryRate)
// 		defer retryTicker.Stop()

// 		timeoutTimer := time.NewTimer(config.Cfg.AckTimeout)
// 		defer timeoutTimer.Stop()

// 		r.SendCh <- msg

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				// new message canceled this send
// 				return
// 			case <-retryTicker.C:
// 				r.SendCh <- msg
// 				fmt.Printf("\n---------------------------ACK-----------------------------------\n")
// 			case ack := <-r.AckIn:
// 				fmt.Println(ack.GetUpdateNr() == msg.GetUpdateNr())
// 				if ack.GetUpdateNr() == msg.GetUpdateNr() {
// 					r.AckResults <- AckResult{ack.GetUpdateNr(), nil}
// 					return
// 				}
// 				// ignore ACKs for older messages or other units
// 			case <-timeoutTimer.C:
// 				r.AckResults <- AckResult{msg.GetUpdateNr(), fmt.Errorf("timeout waiting for ack assignmentSender")}
// 				return
// 			}
// 		}
// 	}()
// }

func (r *GenericSender[A, B]) SendAsyncWithAck(msg A) {
	r.mu.Lock()
	if r.cancelLast != nil {
		r.cancelLast() // cancel previous send
	}

	ctx, cancel := context.WithCancel(context.Background())
	r.cancelLast = cancel
	r.lastUpdateNr = msg.GetUpdateNr()
	r.mu.Unlock()

	ackCh := make(chan B, 1) // per-message ACK channel

	// Optional: register this ackCh with a dispatcher if AckIn is shared
	go func() {
		defer close(ackCh)

		retryTicker := time.NewTicker(config.Cfg.AckRetryRate)
		defer retryTicker.Stop()

		timeoutTimer := time.NewTimer(config.Cfg.AckTimeout)
		defer timeoutTimer.Stop()

		send := func() bool {
			select {
			case r.SendCh <- msg:
				return true
			case <-ctx.Done():
				return false
			}
		}

		if !send() {
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-retryTicker.C:
				send()
			case ack := <-r.AckIn:
				if ack.GetUpdateNr() == msg.GetUpdateNr() {
					select {
					case r.AckResults <- AckResult{ack.GetUpdateNr(), nil}:
					case <-ctx.Done():
					}
					return
				}
				// ignore other ACKs
			case <-timeoutTimer.C:
				select {
				case r.AckResults <- AckResult{msg.GetUpdateNr(), fmt.Errorf("timeout waiting for ack")}:
				case <-ctx.Done():
				}
				return
			}
		}
	}()
}
