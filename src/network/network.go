package network

import (
	"context"
	"elevatorproject/src/config"
	"fmt"
	"sync"
	"time"
)

type OrdersAndStateUpdate struct {
	SourceId string
	UpdateNr int
	OrdersAndState string //Custom type from daniel 
}

type OrdersAndStateAck struct {
	UpdateNr int
}

type AckResult struct {
    UpdateNr int
    Err      error
}

type ReliableSender struct {
    SendCh     chan<- OrdersAndStateUpdate
    AckIn      <-chan OrdersAndStateAck
    AckResults chan AckResult

    cancelLast   context.CancelFunc // cancel previous pending send
    mu           sync.Mutex
    lastUpdateNr int // Must use to prevent confict between different slaves
}


//Sends OrderAndStateUpdates Async to the server. 
func (r *ReliableSender) UpdateAsync(msg OrdersAndStateUpdate, timeout time.Duration) {
    r.mu.Lock()
     // cancel previous send if exists
    if r.cancelLast != nil {
        r.cancelLast()
    }
    ctx, cancel := context.WithCancel(context.Background())
    r.cancelLast = cancel
    r.lastUpdateNr = msg.UpdateNr
    r.mu.Unlock()

    go func() {
        retryTicker := time.NewTicker(config.Cfg.AckRetryRateMs * time.Millisecond)
        defer retryTicker.Stop()

        timeoutTimer := time.NewTimer(timeout)
        defer timeoutTimer.Stop()

        r.SendCh <- msg

        for {
            select {
            case <-ctx.Done():
                // new message canceled this send
                return
            case <-retryTicker.C:
                r.SendCh <- msg
            case ack := <-r.AckIn:
                if ack.UpdateNr == msg.UpdateNr {
                    r.AckResults <- AckResult{UpdateNr: ack.UpdateNr, Err: nil}
                    return
                } 
                // ignore ACKs for older messages ord other units
            case <-timeoutTimer.C:
                r.AckResults <- AckResult{UpdateNr: msg.UpdateNr, Err: fmt.Errorf("timeout waiting for ack")}
                return
            }
        }
    }()
}




