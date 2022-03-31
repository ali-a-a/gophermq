package broker

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/ali-a-a/gophermq/pkg/utils"
	"github.com/google/uuid"
)

var (
	// ErrBadSubject is produced when subject has invalid characters.
	ErrBadSubject = errors.New("invalid subject")
	// ErrMaxPending is produced when overflow occurs.
	ErrMaxPending = errors.New("broker overflow")
	// ErrSubscriberNotFound could be produced in publish.
	ErrSubscriberNotFound = errors.New("subscriber not found")
	// ErrBadID could be produced in fetch.
	ErrBadID = errors.New("bad id")
)

type Broker interface {
	// Publish is responsible to produce data for subject.
	Publish(subject string, data []byte) error
	// Subscribe is responsible to start observing subject.
	Subscribe(subject string) (*subscriber, error)
	// Fetch is responsible to get data from subject.
	Fetch(subject string, id string) ([][]byte, error)
}

// Event is data that is passed through broker.
type Event interface {
	Subject() string
	Data() []byte
	Ack() error
	Error() error
}

// GopherMQ is the Broker implementation.
type GopherMQ struct {
	opts Options

	queue       map[string][][]byte
	pending     map[string]int
	subscribers map[string][]*subscriber
	mutex       sync.Mutex
}

type subscriber struct {
	ID   string
	Subj string
}

// NewGopherMQ returns new GopherMQ.
func NewGopherMQ(opts ...Option) *GopherMQ {
	rand.Seed(time.Now().UnixNano())

	queue := make(map[string][][]byte)
	subscribers := make(map[string][]*subscriber)
	pending := make(map[string]int)

	gm := &GopherMQ{
		queue:       queue,
		subscribers: subscribers,
		pending:     pending,
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
	defer gm.mutex.Unlock()

	if gm.pending[subject] > gm.opts.MaxPending {
		return ErrMaxPending
	}

	subs, ok := gm.subscribers[subject]
	if !ok {
		return ErrSubscriberNotFound
	}

	for _, sub := range subs {
		gm.pending[subject] += 1
		key := utils.SubjectKey(sub.Subj, sub.ID)
		gm.queue[key] = append(gm.queue[key], data)
	}

	return nil
}

func (gm *GopherMQ) Subscribe(subject string) (*subscriber, error) {
	if utils.BadSubject(subject) {
		return nil, ErrBadSubject
	}

	id := uuid.New().String()

	sub := &subscriber{
		ID:   id,
		Subj: subject,
	}

	gm.mutex.Lock()
	gm.subscribers[subject] = append(gm.subscribers[subject], sub)
	gm.mutex.Unlock()

	return sub, nil
}

func (gm *GopherMQ) Fetch(subject string, id string) ([][]byte, error) {
	if utils.BadSubject(subject) {
		return nil, ErrBadSubject
	}

	var target *subscriber

	gm.mutex.Lock()
	subs := gm.subscribers[subject]
	gm.mutex.Unlock()

	for _, sub := range subs {
		if sub.ID == id {
			target = sub
		}
	}

	if target == nil {
		return nil, ErrBadID
	}

	key := utils.SubjectKey(subject, id)

	gm.mutex.Lock()
	data := gm.queue[key]

	gm.pending[subject] -= len(data)
	if gm.pending[subject] < 0 {
		gm.pending[subject] = 0
	}

	delete(gm.queue, key)
	gm.mutex.Unlock()

	return data, nil
}
