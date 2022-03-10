package broker

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGopherMQ_Publish(t *testing.T) {
	gm := NewGopherMQ(MaxPending(1))

	subject := "test"
	data := "data"

	sub, err := gm.Subscribe(subject, func(e Event) error {
		assert.NoError(t, e.Error())
		assert.Equal(t, subject, e.Subject())
		assert.Equal(t, []byte(data), e.Data())

		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, subject, sub.subj)
	assert.NotEmpty(t, sub.id)

	err = gm.Publish(subject, []byte(data))
	assert.NoError(t, err)
	assert.Empty(t, gm.queue[subject])

	time.Sleep(100 * time.Millisecond)

	gm.opts.MaxPending = -1

	err = gm.Publish(subject, []byte(data))
	assert.ErrorIs(t, err, ErrMaxPending)

	gm.opts.MaxPending = 1

	err = gm.Publish("bad subject", []byte(data))
	assert.ErrorIs(t, err, ErrBadSubject)

	_, err = gm.Subscribe(subject, func(e Event) error {
		return errors.New("should be failed")
	})

	assert.NoError(t, err)

	err = gm.Publish(subject, []byte(data))
	assert.Error(t, err)
	assert.NotEqual(t, err, ErrMaxPending)
	assert.NotEmpty(t, gm.queue[subject])

	_ = gm.Publish(subject, []byte(data))
	err = gm.Publish(subject, []byte(data))
	assert.Error(t, err)
	assert.Equal(t, err, ErrMaxPending)
}

func TestGopherMQ_Subscribe(t *testing.T) {
	gm := NewGopherMQ(MaxPending(1))

	subject := "test.a"
	data := "data"

	subject2 := "test.b"
	data2 := "data 2"

	sub, err := gm.Subscribe(subject, func(e Event) error {
		assert.NoError(t, e.Error())
		assert.Equal(t, subject, e.Subject())
		assert.Equal(t, []byte(data), e.Data())

		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, subject, sub.subj)
	assert.NotEmpty(t, sub.id)

	err = gm.Publish(subject, []byte(data))

	assert.NoError(t, err)

	err = gm.Publish(subject, []byte(data))

	assert.NoError(t, err)

	_, err = gm.Subscribe(subject2, func(e Event) error {
		assert.NoError(t, e.Error())
		assert.Equal(t, subject2, e.Subject())
		assert.Equal(t, []byte(data2), e.Data())

		return errors.New("should fail")
	})

	assert.NoError(t, err)

	err = gm.Publish(subject2, []byte(data2))

	assert.Error(t, err)

	var check bool

	_, err = gm.Subscribe(subject2, func(e Event) error {
		assert.NoError(t, e.Error())
		assert.Equal(t, subject2, e.Subject())
		assert.Equal(t, []byte(data2), e.Data())

		check = true

		return nil
	})

	assert.NoError(t, err)

	time.Sleep(100*time.Millisecond)

	assert.True(t, check)
}
