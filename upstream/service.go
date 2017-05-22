package upstream

import (
	"math"
	"math/rand"
	"time"
)

type Service struct {
	tokens       chan struct{}
	mean, stddev time.Duration
	deadline     time.Duration
}

func New(concurrency uint, mean, stddev, timeout time.Duration) *Service {
	s := &Service{
		tokens:   make(chan struct{}, int(concurrency)),
		mean:     mean,
		stddev:   stddev,
		deadline: timeout,
	}

	for i := 0; i < int(concurrency); i++ {
		s.tokens <- struct{}{}
	}

	return s
}

func (s *Service) doWork() <-chan struct{} {
	c := make(chan struct{})

	go func() {
		sleep := time.Duration(math.Abs(rand.NormFloat64())*float64(s.stddev) + float64(s.mean))
		time.Sleep(sleep)
		close(c)
	}()

	return c
}

func (s *Service) timeout() <-chan struct{} {
	c := make(chan struct{})
	go func() {
		time.Sleep(s.deadline)
		close(c)
	}()
	return c
}

func (s *Service) Call() error {
	timeout := s.timeout()

	select {
	case t := <-s.tokens:
		defer func(t struct{}) { s.tokens <- t }(t)
	case <-timeout:
		return QueueTimeout{}
	}

	select {
	case <-s.doWork():
		return nil
	case <-timeout:
		return WorkTimeout{}
	}
}

type QueueTimeout struct{}

func (QT QueueTimeout) Error() string { return "queue timeout expired" }

type WorkTimeout struct{}

func (WT WorkTimeout) Error() string { return "work timeout expired" }
