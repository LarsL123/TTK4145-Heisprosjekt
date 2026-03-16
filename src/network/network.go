package network

import (
	"elevatorproject/src/config"
	"fmt"
	"sync"
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

// func (r *GenericSender[A, B]) SendAsyncWithAck(msg A) {
// 	r.mu.Lock()
// 	if r.cancelLast != nil {
// 		r.cancelLast() // cancel previous send
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	r.cancelLast = cancel
// 	r.lastUpdateNr = msg.GetUpdateNr()
// 	r.mu.Unlock()

// 	ackCh := make(chan B, 1) // per-message ACK channel

// 	// Optional: register this ackCh with a dispatcher if AckIn is shared
// 	go func() {
// 		defer close(ackCh)

// 		retryTicker := time.NewTicker(config.Cfg.AckRetryRate)
// 		defer retryTicker.Stop()

// 		timeoutTimer := time.NewTimer(config.Cfg.AckTimeout)
// 		defer timeoutTimer.Stop()

// 		send := func() bool {
// 			select {
// 			case r.SendCh <- msg:
// 				return true
// 			case <-ctx.Done():
// 				return false
// 			}
// 		}

// 		if !send() {
// 			return
// 		}

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case <-retryTicker.C:
// 				send()
// 			case ack := <-r.AckIn:
// 				if ack.GetUpdateNr() == msg.GetUpdateNr() {
// 					select {
// 					case r.AckResults <- AckResult{ack.GetUpdateNr(), nil}:
// 					case <-ctx.Done():
// 					}
// 					return
// 				}
// 				// ignore other ACKs
// 			case <-timeoutTimer.C:
// 				select {
// 				case r.AckResults <- AckResult{msg.GetUpdateNr(), fmt.Errorf("timeout waiting for ack")}:
// 				case <-ctx.Done():
// 				}
// 				return
// 			}
// 		}
// 	}()
// }
type GenericSender[SenderType NetMessage, ReciverType NetMessage] struct {
    SendCh     chan<- SenderType
    AckIn      <-chan ReciverType
    AckResults chan AckResult

    mu           sync.Mutex
    lastUpdateNr int
    ackMap       sync.Map // map[int]chan ReciverType
}

// Dispatcher forwards ACKs to the correct channel
func (r *GenericSender[A, B]) StartAckDispatcher() {
    go func() {
        for ack := range r.AckIn {
            updateNr := ack.GetUpdateNr()
            if chInterface, ok := r.ackMap.Load(updateNr); ok {
                ch := chInterface.(chan B)
                select {
                case ch <- ack:
                    // delivered
                default:
                    // goroutine already finished (timeout or done)
                }
            }
        }
    }()
}

func (r *GenericSender[A, B]) SendAsyncWithAck(msg A) {
    updateNr := msg.GetUpdateNr()

    // Each message gets its own ACK channel
    ackCh := make(chan B, 1)
    r.ackMap.Store(updateNr, ackCh)
    defer func() {
        // Clean up after ACK or timeout
        r.ackMap.Delete(updateNr)
    }()

    go func() {
        retryTicker := time.NewTicker(config.Cfg.AckRetryRate)
        defer retryTicker.Stop()

        timeoutTimer := time.NewTimer(config.Cfg.AckTimeout)
        defer timeoutTimer.Stop()

        r.SendCh <- msg // initial send

        for {
            select {
            case <-retryTicker.C:
                r.SendCh <- msg // retry
                fmt.Printf("\n---------------------------ACK RETRY-----------------------------------\n")
            case ack := <-ackCh:
                // ACK received
                r.AckResults <- AckResult{ack.GetUpdateNr(), nil}
                return
            case <-timeoutTimer.C:
                r.AckResults <- AckResult{updateNr, fmt.Errorf("timeout waiting for ack")}
                return
            }
        }
    }()
}
// type GenericSender [SenderType NetMessage, ReciverType NetMessage] struct {
//     SendCh     chan<- SenderType
//     AckIn      <-chan ReciverType
//     AckResults chan AckResult

//     cancelLast   context.CancelFunc // cancel previous pending send
//     mu           sync.Mutex 
//     lastUpdateNr int // Must use to prevent confict between different slaves
// }

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


// func (r *GenericSender[A, B]) SendAsyncWithAck(msg A) {
//     r.mu.Lock()
//     if r.cancelLast != nil {
//         r.cancelLast() // cancel previous send
//     }

//     ctx, cancel := context.WithCancel(context.Background())
//     r.cancelLast = cancel
//     r.lastUpdateNr = msg.GetUpdateNr()
//     r.mu.Unlock()

//     ackCh := make(chan B, 1) // per-message ACK channel

//     // Optional: register this ackCh with a dispatcher if AckIn is shared
//     go func() {
//         defer close(ackCh)

//         retryTicker := time.NewTicker(config.Cfg.AckRetryRate)
//         defer retryTicker.Stop()

//         timeoutTimer := time.NewTimer(config.Cfg.AckTimeout)
//         defer timeoutTimer.Stop()

//         send := func() bool {
//             select {
//             case r.SendCh <- msg:
//                 return true
//             case <-ctx.Done():
//                 return false
//             }
//         }

//         if !send() {
//             return
//         }

//         for {
//             select {
//             case <-ctx.Done():
//                 return
//             case <-retryTicker.C:
//                 send()
//             case ack := <-r.AckIn:
//                 if ack.GetUpdateNr() == msg.GetUpdateNr() {
//                     select {
//                     case r.AckResults <- AckResult{ack.GetUpdateNr(), nil}:
//                     case <-ctx.Done():
//                     }
//                     return
//                 }
//                 // ignore other ACKs
//             case <-timeoutTimer.C:
//                 select {
//                 case r.AckResults <- AckResult{msg.GetUpdateNr(), fmt.Errorf("timeout waiting for ack")}:
//                 case <-ctx.Done():
//                 }
//                 return
//             }
//         }
//     }()
// }


// AckResult represents the result of a sent message
type AckResult struct {
	UpdateNr int
	Err      error
}

// AckListener manages ACKs for sent messages
type AckListener struct {
	pending     map[int]chan AckResult
	mutex       sync.Mutex
	timeout     time.Duration
	incomingAck chan int // receives UpdateNr from network
}

// NewAckListener creates a new AckListener
func NewAckListener(timeout time.Duration) *AckListener {
	al := &AckListener{
		pending:     make(map[int]chan AckResult),
		timeout:     timeout,
		incomingAck: make(chan int, 100),
	}
	go al.run()
	return al
}

// WaitAck registers a message and waits for ACK or timeout
func (al *AckListener) WaitAck(updateNr int) AckResult {
	al.mutex.Lock()
	ch := make(chan AckResult, 1)
	al.pending[updateNr] = ch
	al.mutex.Unlock()

	select {
	case res := <-ch:
		return res
	case <-time.After(al.timeout):
		al.removePending(updateNr)
		return AckResult{UpdateNr: updateNr, Err: fmt.Errorf("ACK timeout")}
	}
}

// HandleAck should be called when a network ACK arrives
func (al *AckListener) HandleAck(updateNr int) {
	al.mutex.Lock()
	defer al.mutex.Unlock()
	if ch, ok := al.pending[updateNr]; ok {
		ch <- AckResult{UpdateNr: updateNr, Err: nil}
		delete(al.pending, updateNr)
	}
}

// internal loop to process incoming ACKs
func (al *AckListener) run() {
	for updateNr := range al.incomingAck {
		al.HandleAck(updateNr)
	}
}

// helper to remove pending on timeout
func (al *AckListener) removePending(updateNr int) {
	al.mutex.Lock()
	defer al.mutex.Unlock()
	delete(al.pending, updateNr)
}

// External function to push ACK from network
func (al *AckListener) PushAck(updateNr int) {
	al.incomingAck <- updateNr
}