package pubsub

import (
	"context"
	"errors"
	"sync"
	"time"
)

// PubSub contains and manage the map of topics -> subscribers
type PubSub[T comparable, P any] struct {
	sync.Mutex
	m map[T][]*Sub[T, P]
}

func NewPubSub[T comparable, P any]() *PubSub[T, P] {
	ps := PubSub[T, P]{}
	ps.m = make(map[T][]*Sub[T, P])
	return &ps
}

func (p *PubSub[T, P]) getSubscribers(topic T) []*Sub[T, P] {
	p.Lock()
	defer p.Unlock()
	return p.m[topic]
}

func (p *PubSub[T, P]) addSubscriber(s *Sub[T, P]) {
	p.Lock()
	for _, topic := range s.topics {
		p.m[topic] = append(p.m[topic], s)
	}
	p.Unlock()
}

func (p *PubSub[T, P]) removeSubscriber(s *Sub[T, P]) {
	p.Lock()
	for _, topic := range s.topics {
		for i, subscriber := range p.m[topic] {
			if subscriber == s {
				p.m[topic] = append(p.m[topic][:i], p.m[topic][i+1:]...)
				break
			}
		}
	}
	p.Unlock()
}

// Subscribe is an alias for NewSub
func (p *PubSub[T, P]) Subscribe(topics []T) *Sub[T, P] {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Sub[T, P]{topics: topics, ch: make(chan Payload[T, P], 10), ctx: ctx, cancel: cancel, p: p}
	p.addSubscriber(s)
	return s
}

// Pub shortcut for publish which ignore the error
func (p *PubSub[T, P]) Pub(topic T, msg P) {
	for _, s := range p.getSubscribers(topic) {
		s.publish(Payload[T, P]{topic, msg})
	}
}

type Payload[T comparable, P any] struct {
	Topic T
	Msg   P
}

// ErrTimeout error returned when timeout occurs
var ErrTimeout = errors.New("timeout")

// ErrCancelled error returned when context is cancelled
var ErrCancelled = errors.New("cancelled")

// Sub subscriber will receive messages published on a Topic in his ch
type Sub[T comparable, P any] struct {
	topics []T                // Topics subscribed to
	ch     chan Payload[T, P] // Receives messages in this channel
	ctx    context.Context
	cancel context.CancelFunc
	p      *PubSub[T, P]
}

// ReceiveTimeout2 returns a message received on the channel or timeout
func (s *Sub[T, P]) ReceiveTimeout2(timeout time.Duration, c1 <-chan struct{}) (topic T, msg P, err error) {
	select {
	case p := <-s.ch:
		return p.Topic, p.Msg, nil
	case <-time.After(timeout):
		return topic, msg, ErrTimeout
	case <-c1:
		return topic, msg, ErrCancelled
	case <-s.ctx.Done():
		return topic, msg, ErrCancelled
	}
}

// ReceiveTimeout returns a message received on the channel or timeout
func (s *Sub[T, P]) ReceiveTimeout(timeout time.Duration) (topic T, msg P, err error) {
	c1 := make(chan struct{})
	return s.ReceiveTimeout2(timeout, c1)
}

// Receive returns a message
func (s *Sub[T, P]) Receive() (topic T, msg P, err error) {
	var res P
	select {
	case p := <-s.ch:
		return p.Topic, p.Msg, nil
	case <-s.ctx.Done():
		return topic, res, ErrCancelled
	}
}

// ReceiveCh returns a message
func (s *Sub[T, P]) ReceiveCh() <-chan Payload[T, P] {
	return s.ch
}

// Close will remove the subscriber from the Topic subscribers
func (s *Sub[T, P]) Close() {
	s.cancel()
	s.p.removeSubscriber(s)
}

// publish a message to the subscriber channel
func (s *Sub[T, P]) publish(p Payload[T, P]) {
	select {
	case s.ch <- p:
	default:
	}
}
