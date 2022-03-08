package broker

import (
	"context"
	"errors"
	"github.com/ali-a-a/gophermq/pkg/utils"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// ErrBadSubject is produced when subject has invalid characters.
	ErrBadSubject = errors.New("invalid subject")
)

type Broker interface {
	// Publish is responsible to produce data for subject.
	Publish(subject string, data []byte) error
}

// GopherMQ is the Broker implementation.
type GopherMQ struct {
	queue            map[string][][]byte
	numOfSubscribers map[string]*uint32
	timeout          time.Duration
	mutex            sync.Mutex
}

// NewGopherMQ returns new GopherMQ.
func NewGopherMQ() *GopherMQ {
	queue := make(map[string][][]byte)

	return &GopherMQ{queue: queue}
}

func (gm *GopherMQ) Publish(subject string, data []byte) error {
	if utils.BadSubject(subject) {
		return ErrBadSubject
	}

	gm.mutex.Lock()
	gm.queue[subject] = append(gm.queue[subject], data)
	gm.mutex.Unlock()

	atomic.AddUint32(gm.numOfSubscribers[subject], 1)

	ctx, cancel := context.WithTimeout(context.Background(), gm.timeout)

	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:

		}
	}

	return nil
}
