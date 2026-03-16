package network

import (
	"context"
	"elevatorproject/src/config"
	"fmt"
	"time"
)

func (r *GenericSender[A, B]) SendAsyncWithAck(msg A) {
	r.mu.Lock()
	// cancel previous send if exists
	if r.cancelLast != nil {
		r.cancelLast()
	}
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelLast = cancel
	r.lastUpdateNr = msg.GetUpdateNr()
	r.mu.Unlock()

	go func() {
		retryTicker := time.NewTicker(config.Cfg.AckRetryRate)
		defer retryTicker.Stop()

		timeoutTimer := time.NewTimer(config.Cfg.AckTimeout)
		defer timeoutTimer.Stop()

		r.SendCh <- msg

		for {
			select {
			case <-ctx.Done():
				// new message canceled this send
				return
			case <-retryTicker.C:
				r.SendCh <- msg
				fmt.Printf("\n---------------------------ACK-----------------------------------\n")
			case ack := <-r.AckIn:
				fmt.Println(ack.GetUpdateNr() == msg.GetUpdateNr())
				if ack.GetUpdateNr() == msg.GetUpdateNr() {
					r.AckResults <- AckResult{ack.GetUpdateNr(), nil}
					return
				}
				// ignore ACKs for older messages or other units
			case <-timeoutTimer.C:
				r.AckResults <- AckResult{msg.GetUpdateNr(), fmt.Errorf("timeout waiting for ack assignmentSender")}
				return
			}
		}
	}()
}

//Sends OrderAndStateUpdates Async to the server. Sending a new update will cancel the previous one if its still trying to send.
// func (r *AssignmentSender) UpdateAsync(msg AssignmentsAndOrders) {
//     r.mu.Lock()
//      // cancel previous send if exists
//     if r.cancelLast != nil {
//         r.cancelLast()
//     }
//     ctx, cancel := context.WithCancel(context.Background())
//     r.cancelLast = cancel
//     r.lastUpdateNr = msg.UpdateNr
//     r.mu.Unlock()

//     go func() {
//         retryTicker := time.NewTicker(config.Cfg.AckRetryRate)
//         defer retryTicker.Stop()

//         timeoutTimer := time.NewTimer(config.Cfg.AckTimeout)
//         defer timeoutTimer.Stop()

//         r.SendCh <- msg

//         for {
//             select {
//             case <-ctx.Done():
//                 // new message canceled this send
//                 return
//             case <-retryTicker.C:
//                 r.SendCh <- msg
//             case ack := <-r.AckIn:
//                 if ack.UpdateNr == msg.UpdateNr {
//                     r.AckResults <- AckResult{UpdateNr: ack.UpdateNr, Err: nil}
//                     return
//                 }
//                 // ignore ACKs for older messages ord other units
//             case <-timeoutTimer.C:
//                 r.AckResults <- AckResult{UpdateNr: msg.UpdateNr, Err: fmt.Errorf("timeout waiting for ack assignmentSender")}
//                 return
//             }
//         }
//     }()
// }

//Sends OrderAndStateUpdates Async to the server. Seding a new update will cancel the previous one if its still trying to send.
// func (r *OrderSender) UpdateAsync(msg OrdersAndStateUpdate) {
//     r.mu.Lock()
//      // cancel previous send if exists
//     if r.cancelLast != nil {
//         r.cancelLast()
//     }
//     ctx, cancel := context.WithCancel(context.Background())
//     r.cancelLast = cancel
//     r.lastUpdateNr = msg.UpdateNr
//     r.mu.Unlock()

//     go func() {
//         retryTicker := time.NewTicker(config.Cfg.AckRetryRate)
//         defer retryTicker.Stop()

//         timeoutTimer := time.NewTimer(config.Cfg.AckTimeout)
//         defer timeoutTimer.Stop()

//         r.SendCh <- msg

//         for {
//             select {
//             case <-ctx.Done():
//                 // new message canceled this send
//                 return
//             case <-retryTicker.C:
//                 r.SendCh <- msg
//             case ack := <-r.AckIn:
//                 if ack.UpdateNr == msg.UpdateNr {
//                     r.AckResults <- AckResult{UpdateNr: ack.UpdateNr, Err: nil}
//                     return
//                 }
//                 // ignore ACKs for older messages ord other units
//             case <-timeoutTimer.C:
//                 r.AckResults <- AckResult{UpdateNr: msg.UpdateNr, Err: fmt.Errorf("timeout waiting for ack orderSender")}
//                 return
//             }
//         }
//     }()
// }
