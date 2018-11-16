package app

type Heartbeater struct {
	beatFunc func() error
	errCh    chan error
	quitCh   chan struct{}
	doneCh   chan struct{}
}

func NewHeartbeater(beatFunc func() error) (*Heartbeater, error) {
	return &Heartbeater{
		beatFunc: beatFunc,
		errCh:    make(chan error, 1),
		quitCh:   make(chan struct{}, 0),
		doneCh:   make(chan struct{}, 0),
	}, nil
}

func (h *Heartbeater) Beat() (bool, error) {
	go func() {
		h.errCh <- h.beatFunc()
	}()

	select {
	case err := <-h.errCh:
		return true, err
	case <-h.quitCh:
		close(h.doneCh)
		return false, nil
	}
}

func (h *Heartbeater) Quit() {
	close(h.quitCh)
	<-h.doneCh
}
