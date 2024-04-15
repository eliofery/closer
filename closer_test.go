package closer

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

// CloserMock implements Closer.
type CloserMock struct {
	wg     sync.WaitGroup
	isStop bool
}

// NewMock creates a new instance of CloserMock.
func NewMock(sigs ...os.Signal) *CloserMock {
	c := &CloserMock{}

	if len(sigs) == 0 {
		sigs = DefaultSignals
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, sigs...)

		time.Sleep(1 * time.Second)
		signal.Stop(ch)
	}()

	return c
}

// Wait waits for all functions to be closed.
func (c *CloserMock) Wait() {
	c.wg.Wait()
	c.isStop = true
}

// TestNewAndWait tests the New function and waits for a signal.
func TestNewAndWait(t *testing.T) {
	clr := NewMock(DefaultSignals...)
	clr.Wait()

	assert.Equal(t, clr.isStop, true)
}

// TestAddAndClose tests the Add and Close functions.
func TestAddAndClose(t *testing.T) {
	clr := New()

	clr.Add(func() error {
		return nil
	})

	clr.Add(func() error {
		return errors.New("error")
	})

	err := clr.Close()

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "error")
	assert.Equal(t, nil, clr.Close())
}
