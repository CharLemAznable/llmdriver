package llmdriver

import (
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/os/gmutex"
)

func NewDefaultOutputStream() *DefaultOutputStream {
	return &DefaultOutputStream{
		ch:     make(chan OutputEvent),
		mutex:  &gmutex.Mutex{},
		closed: gtype.NewBool(),
	}
}

type DefaultOutputStream struct {
	ch     chan OutputEvent
	mutex  *gmutex.Mutex
	err    error
	closed *gtype.Bool
}

func (s *DefaultOutputStream) Push(event OutputEvent) {
	s.ch <- event
}

func (s *DefaultOutputStream) Close(err error) {
	if !s.closed.Cas(false, true) {
		return // already closed
	}
	s.mutex.LockFunc(func() {
		s.err = err
		close(s.ch)
	})
}

func (s *DefaultOutputStream) Event() <-chan OutputEvent {
	return s.ch
}

func (s *DefaultOutputStream) Err() (err error) {
	s.mutex.LockFunc(func() {
		err = s.err
	})
	return
}

func (s *DefaultOutputStream) Drain() {
	for range s.ch {
		// drain channel
	}
}
