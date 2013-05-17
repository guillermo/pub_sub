// License: BSD

// Publish Subscribe with go exchanges
//
// This library provie the pub/sub pattern to the go routines
package pub_sub

import (
	"sync"
)

type Exchange struct {
	mu          sync.Mutex
	subscribers []*Subscription
}

type Subscription struct {
	C        chan interface{}
	exchange *Exchange
}

// Create a new Publish/Subscribe exchange
func NewPubSub() *Exchange {
	exchange := &Exchange{subscribers: make([]*Subscription, 0)}
	return exchange
}

func (e *Exchange) unsubscribe(s *Subscription) {
}

func (c *Exchange) subscribe(s *Subscription) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.subscribers = append(c.subscribers, s)
}

// Subscribe to exchange
//
//     subscription := exchange.Subscribe()
//     msg := <- sbuscription.C
//
func (e *Exchange) Subscribe() *Subscription {
	subscription := &Subscription{make(chan interface{}), e}
	e.subscribe(subscription)
	return subscription
}

// Publish a message into the exchange (Broadcast)
// It will go to all the subscriptions and send individually the message
// It will block until all the subscribers recive the message
// You may want to launch this in its independent gorutine
//
//     go exchange.Publish("msg")
//
func (e *Exchange) Publish(data interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, subscriber := range e.subscribers {
		subscriber.C <- data
	}
}

// Unbscribe the subscription
// This step is necessary to skip memory leaks
func (s *Subscription) Unsubscribe() {
	s.exchange.mu.Lock()
	defer s.exchange.mu.Unlock()
	for i, subscriber := range s.exchange.subscribers {
		if subscriber == s {
			n := i + 1
			left := s.exchange.subscribers[:i]
			right := s.exchange.subscribers[n:]
			s.exchange.subscribers = append(left, right...)
		}
	}
}
