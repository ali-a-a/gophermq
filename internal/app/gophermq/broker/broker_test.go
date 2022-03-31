package broker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGopherMQ_Publish(t *testing.T) {
	gm := NewGopherMQ(MaxPending(1))

	subject := "test"
	data := "data"

	sub, err := gm.Subscribe(subject)

	assert.NoError(t, err)
	assert.Equal(t, subject, sub.Subj)
	assert.NotEmpty(t, sub.ID)

	err = gm.Publish(sub.Subj, []byte(data))
	assert.NoError(t, err)

	d, err := gm.Fetch(sub.Subj, sub.ID)
	assert.NoError(t, err)
	assert.Equal(t, []byte(data), d[0])

	sub2, err := gm.Subscribe(subject)

	assert.NoError(t, err)
	assert.Equal(t, subject, sub2.Subj)
	assert.NotEmpty(t, sub2.ID)

	err = gm.Publish(sub.Subj, []byte(data))
	assert.NoError(t, err)

	d1, err := gm.Fetch(sub.Subj, sub.ID)
	assert.NoError(t, err)
	assert.Equal(t, []byte(data), d1[0])

	d2, err := gm.Fetch(sub2.Subj, sub2.ID)
	assert.NoError(t, err)
	assert.Equal(t, []byte(data), d2[0])

	d3, err := gm.Fetch(sub2.Subj, sub2.ID)
	assert.NoError(t, err)
	assert.Empty(t, d3)

	subject2 := "test2"

	sub3, err := gm.Subscribe(subject2)

	assert.NoError(t, err)
	assert.Equal(t, subject2, sub3.Subj)
	assert.NotEmpty(t, sub3.ID)

	err = gm.Publish(sub3.Subj, []byte(data))
	assert.NoError(t, err)

	d4, err := gm.Fetch(sub3.Subj, sub3.ID)
	assert.NoError(t, err)
	assert.Equal(t, []byte(data), d4[0])

	err = gm.Publish("bad subject", []byte(data))
	assert.Equal(t, ErrBadSubject, err)

	err = gm.Publish("subscriber.not.found", []byte(data))
	assert.Equal(t, ErrSubscriberNotFound, err)

	gm.opts.MaxPending = -1

	err = gm.Publish(subject, []byte(data))
	assert.Equal(t, ErrMaxPending, err)
}

func TestGopherMQ_Subscribe(t *testing.T) {
	gm := NewGopherMQ(MaxPending(1))

	subject := "test.a"
	subject2 := "test.b"
	subject3 := "test.c"
	subject32 := "test.c"

	sub, err := gm.Subscribe(subject)
	sub2, err2 := gm.Subscribe(subject2)
	sub3, err3 := gm.Subscribe(subject3)
	sub32, err4 := gm.Subscribe(subject32)

	assert.NoError(t, err)
	assert.Equal(t, subject, sub.Subj)
	assert.NotEmpty(t, sub.ID)

	assert.NoError(t, err2)
	assert.Equal(t, subject2, sub2.Subj)
	assert.NotEmpty(t, sub2.ID)

	assert.NoError(t, err3)
	assert.Equal(t, subject3, sub3.Subj)
	assert.NotEmpty(t, sub3.ID)

	assert.NoError(t, err4)
	assert.Equal(t, subject32, sub32.Subj)
	assert.NotEmpty(t, sub32.ID)

	assert.NotEqual(t, sub.ID, sub2.ID)
	assert.NotEqual(t, sub.ID, sub3.ID)
	assert.NotEqual(t, sub2.ID, sub3.ID)
	assert.NotEqual(t, sub3.ID, sub32.ID)
}

func TestGopherMQ_Fetch(t *testing.T) {
	gm := NewGopherMQ(MaxPending(1))

	subject := "test"
	data := "data"

	sub, err := gm.Subscribe(subject)

	assert.NoError(t, err)
	assert.Equal(t, subject, sub.Subj)
	assert.NotEmpty(t, sub.ID)

	err = gm.Publish(sub.Subj, []byte(data))
	assert.NoError(t, err)

	d, err := gm.Fetch(sub.Subj, sub.ID)
	assert.NoError(t, err)
	assert.Equal(t, []byte(data), d[0])

	d, err = gm.Fetch(sub.Subj, "bad id")
	assert.Empty(t, d)
	assert.Equal(t, ErrBadID, err)
}
