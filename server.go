package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

type state struct {
	ratio   int64
	closing uint32
}

func (s *state) needClose() bool {
	return atomic.LoadUint32(&s.closing) == 1 || rand.Int63n(100) < s.ratio
}

func (s *state) setClosing() {
	atomic.StoreUint32(&s.closing, 1)
}

type Server struct {
	Server         *http.Server
	WaitBeforeStop time.Duration
}

func (s *Server) Run(ctx0 context.Context) error {
	st := &state{ratio: 5, closing: 0}
	h := s.Server.Handler
	s.Server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if st.needClose() {
			w.Header().Set("Connection", "close")
		}
		h.ServeHTTP(w, r)
	})

	ctx1, cancel1 := context.WithCancel(ctx0)
	defer cancel1()
	shutdown := make(chan error, 1)
	go func() {
		defer close(shutdown)
		<-ctx1.Done()
		select {
		case <-ctx0.Done():
		default:
			return
		}
		if wait := s.WaitBeforeStop; wait != 0 {
			st.setClosing()
			log.Printf("wait %s before stopping...", wait.String())
			time.Sleep(wait)
		}
		ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel2()
		if err := s.Server.Shutdown(ctx2); err != nil {
			shutdown <- err
		}
	}()

	err1 := s.Server.ListenAndServe()
	cancel1()
	err2 := <-shutdown
	if err1 != nil && !errors.Is(err1, http.ErrServerClosed) {
		return err1
	}
	return err2
}
