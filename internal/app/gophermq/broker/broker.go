package broker

import (
	"errors"
	"github.com/ali-a-a/gophermq/pkg/utils"
	"github.com/google/uuid"
	"math/rand"
	"sync"
	"time"
)

var (
	// ErrBadSubject is produced when subject has invalid characters.
	ErrBadSubject = errors.New("invalid subject")
	// ErrMaxPending is produced when overflow occurs.
	ErrMaxPending = errors.New("broker overflow")
)

type Broker interface {
	// Publish is responsible to produce data for subject.
	Publish(subject string, data []byte) error
	// Subscribe is responsible to consume data from subject.
	Subscribe(subject string, handler Handler) (*subscriber, error)
}

// Event is data that is passed through broker.
type Event interface {
	Subject() string
	Data() []byte
	Ack() error
	Error() error
}

// Handler is used to process messages via a subscription of a subject.
type Handler func(Event) error

// GopherMQ is the Broker implementation.
type GopherMQ struct {
	opts Options

	queue       map[string][][]byte
	subscribers map[string][]*subscriber
	mutex       sync.Mutex
}

type subscriber struct {
	id      string
	subj    string
	handler Handler
}

type event struct {
	subj string
	err  error
	data []byte
}

func (m *event) Subject() string {
	return m.subj
}

func (m *event) Data() []byte {
	return m.data
}

func (m *event) Ack() error {
	return nil
}

func (m *event) Error() error {
	return m.err
}

// NewGopherMQ returns new GopherMQ.
func NewGopherMQ(opts ...Option) *GopherMQ {
	rand.Seed(time.Now().UnixNano())

	queue := make(map[string][][]byte)
	subscribers := make(map[string][]*subscriber)

	gm := &GopherMQ{
		queue:       queue,
		subscribers: subscribers,
	}

	for _, o := range opts {
		o(&gm.opts)
	}

	return gm
}

func (gm *GopherMQ) Publish(subject string, data []byte) error {
	if utils.BadSubject(subject) {
		return ErrBadSubject
	}

	gm.mutex.Lock()

	if len(gm.queue[subject]) > gm.opts.MaxPending {
		gm.mutex.Unlock()

		return ErrMaxPending
	}

	gm.queue[subject] = append(gm.queue[subject], data)

	subs, ok := gm.subscribers[subject]

	gm.mutex.Unlock()
	if !ok {
		return nil
	}

	e := &event{
		subj: subject,
		data: data,
	}

	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	for _, sub := range subs {
		if err := sub.handler(e); err != nil {
			return err
		}
	}

	delete(gm.queue, subject)

	return nil
}

func (gm *GopherMQ) Subscribe(subject string, handler Handler) (*subscriber, error) {
	if utils.BadSubject(subject) {
		return nil, ErrBadSubject
	}

	sub := &subscriber{
		id:      uuid.New().String(),
		subj:    subject,
		handler: handler,
	}

	gm.mutex.Lock()
	gm.subscribers[subject] = append(gm.subscribers[subject], sub)
	gm.mutex.Unlock()

	go func() {
		for _, msg := range gm.queue[subject] {
			_ = handler(&event{
				subj: subject,
				data: msg,
			})
		}

		delete(gm.queue, subject)
	}()

	return sub, nil
}
