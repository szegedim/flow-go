package module

import (
	"errors"
	"sync"

	"github.com/onflow/flow-go/module/irrecoverable"
)

// WARNING: The semantics of this interface will be changing in the near future, with
// startup / shutdown capabilities being delegated to the Startable interface instead.
// For more details, see [FLIP 1167](https://github.com/onflow/flow-go/pull/1167)
//
// ReadyDoneAware provides an easy interface to wait for module startup and shutdown.
// Modules that implement this interface only support a single start-stop cycle, and
// will not restart if Ready() is called again after shutdown has already commenced.
type ReadyDoneAware interface {
	// Ready commences startup of the module, and returns a ready channel that is closed once
	// startup has completed. Note that the ready channel may never close if errors are
	// encountered during startup.
	// If shutdown has already commenced before this method is called for the first time,
	// startup will not be performed and the returned channel will also never close.
	// This should be an idempotent method.
	Ready() <-chan struct{}

	// Done commences shutdown of the module, and returns a done channel that is closed once
	// shutdown has completed. Note that the done channel should be closed even if errors are
	// encountered during shutdown.
	// This should be an idempotent method.
	Done() <-chan struct{}
}

type NoopReadDoneAware struct{}

func (n *NoopReadDoneAware) Ready() <-chan struct{} {
	ready := make(chan struct{})
	defer close(ready)
	return ready
}

func (n *NoopReadDoneAware) Done() <-chan struct{} {
	done := make(chan struct{})
	defer close(done)
	return done
}

var ErrMultipleStartup = errors.New("component may only be started once")

// Startable provides an interface to start a component. Once started, the component
// can be stopped by cancelling the given context.
type Startable interface {
	// Start starts the component. Any errors encountered during startup should be returned
	// directly, whereas irrecoverable errors encountered while the component is running
	// should be thrown with the given SignalerContext.
	// This method should only be called once, and subsequent calls should return ErrMultipleStartup.
	Start(irrecoverable.SignalerContext) error
}

// AllReady calls Ready on all input components and returns a channel that is
// closed when all input components are ready.
func AllReady(components ...ReadyDoneAware) <-chan struct{} {
	ready := make(chan struct{})
	var wg sync.WaitGroup

	for _, component := range components {
		wg.Add(1)
		go func(c ReadyDoneAware) {
			<-c.Ready()
			wg.Done()
		}(component)
	}

	go func() {
		wg.Wait()
		close(ready)
	}()

	return ready
}

// AllDone calls Done on all input components and returns a channel that is
// closed when all input components are done.
func AllDone(components ...ReadyDoneAware) <-chan struct{} {
	done := make(chan struct{})
	var wg sync.WaitGroup

	for _, component := range components {
		wg.Add(1)
		go func(c ReadyDoneAware) {
			<-c.Done()
			wg.Done()
		}(component)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	return done
}
