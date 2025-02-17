package services

import "context"

type (
	// Service represents a long running service inside the
	// Application.
	//
	// Typically a Service will leverage utils.StartStopOnce to implement these
	// calls in a safe manner.
	//
	// Template
	//
	// Mockable Foo service with a run loop
	//  //go:generate mockery --name Foo --output ../internal/mocks/ --case=underscore
	//  type (
	//  	// Expose a public interface so we can mock the service.
	//  	Foo interface {
	//  		service.Service
	//
	//  		// ...
	//  	}
	//
	//  	foo struct {
	//  		// ...
	//
	//  		stop chan struct{}
	//  		done chan struct{}
	//
	//  		utils.StartStopOnce
	//  	}
	//  )
	//
	//  var _ Foo = (*foo)(nil)
	//
	//  func NewFoo() Foo {
	//  	f := &foo{
	//  		// ...
	//  	}
	//
	//  	return f
	//  }
	//
	//  func (f *foo) Start() error {
	//  	return f.StartOnce("Foo", func() error {
	//  		go f.run()
	//
	//  		return nil
	//  	})
	//  }
	//
	//  func (f *foo) Close() error {
	//  	return f.StopOnce("Foo", func() error {
	//  		// trigger goroutine cleanup
	//  		close(f.stop)
	//  		// wait for cleanup to complete
	//  		<-f.done
	//  		return nil
	//  	})
	//  }
	//
	//  func (f *foo) run() {
	//  	// signal cleanup completion
	//  	defer close(f.done)
	//
	//  	for {
	//  		select {
	//  		// ...
	//  		case <-f.stop:
	//  			// stop the routine
	//  			return
	//  		}
	//  	}
	//
	//  }
	Service interface {
		// Start the service.
		Start() error
		// Close stops the Service.
		// Invariants: Usually after this call the Service cannot be started
		// again, you need to build a new Service to do so.
		Close() error

		Checkable
	}

	// ServiceCtx is the same Service interface, but Start function receives a context.
	// This is needed for services that make HTTP calls or DB queries in Start.
	ServiceCtx interface {
		// Start the service. Must quit immediately if the context is cancelled.
		// The given context applies to Start function only and must not be retained.
		Start(context.Context) error
		// Close stops the Service.
		// Invariants: Usually after this call the Service cannot be started
		// again, you need to build a new Service to do so.
		Close() error

		Checkable
	}
)

type adapter struct {
	service Service
}

// NewServiceCtx creates an adapter instance for the given Service.
func NewServiceCtx(service Service) ServiceCtx {
	return &adapter{
		service,
	}
}

// Start forwards the call to the underlying service.Start().
// Context is not used in this case.
func (a adapter) Start(context.Context) error {
	return a.service.Start()
}

// Close forwards the call to the underlying service.Close().
func (a adapter) Close() error {
	return a.service.Close()
}

// Ready forwards the call to the underlying service.Ready().
func (a adapter) Ready() error {
	return a.service.Ready()
}

// Healthy forwards the call to the underlying service.Healthy().
func (a adapter) Healthy() error {
	return a.service.Healthy()
}
