package closer

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	// DefaultSignals is the default list of signals.
	DefaultSignals = []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	}
)

// handler is a function that can be added to Closer.
type handler func() error

// Closer is an interface for adding and closing functions.
type Closer interface {
	// Wait waits for all functions to be closed.
	Wait()

	// Add adds a function to be closed.
	Add(fn handler)

	// Close closes all functions.
	Close() error
}

// closer implements Closer.
type closer struct {
	wg   sync.WaitGroup
	mu   sync.Mutex
	once sync.Once
	fns  []handler
}

// New creates a new instance of Closer that waits for the signals passed as arguments.
// If no signals are specified, default values (defaultSignals) are used.
// The function starts a goroutine that waits for one of the specified signals to be
// received and exits when it occurs.
// Upon receiving a signal, the goroutine stops monitoring signals.
func New(sigs ...os.Signal) Closer {
	c := &closer{}

	if len(sigs) == 0 {
		sigs = DefaultSignals
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, sigs...)

		select {
		case <-ch:
			signal.Stop(ch)
		}
	}()

	return c
}

// Wait waits for all functions to be closed.
func (c *closer) Wait() {
	c.wg.Wait()
}

// Add adds a function to be closed.
func (c *closer) Add(fn handler) {
	c.mu.Lock()
	c.fns = append(c.fns, fn)
	c.mu.Unlock()
}

// Close closes all functions.
func (c *closer) Close() error {
	var errs []error

	c.once.Do(func() {
		var wg sync.WaitGroup
		errsCh := make(chan error, len(c.fns))

		for _, fn := range c.fns {
			wg.Add(1)

			go func(fn handler) {
				defer wg.Done()

				if err := fn(); err != nil {
					errsCh <- err
				}
			}(fn)
		}

		go func() {
			wg.Wait()
			close(errsCh)
		}()

		for {
			err, ok := <-errsCh
			if !ok {
				break
			}

			errs = append(errs, err)
		}
	})

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
