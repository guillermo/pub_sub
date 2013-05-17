// License: BSD

// Publish Subscribe with go channels
//
// This library provie the pub/sub pattern to the go routines
package pub_sub

import (
	"sync"
)

type Exchange struct {
	sync.Mutex
	subscriptors [](*Subscription)
}

type Subscription struct {
	C       chan (interface{})
	channel *Exchange
}

// Create a new Publish/Subscribe Channel
func NewPubSub() *Exchange {
	channel := &Exchange{subscriptors: make([](*Subscription), 0)}
	return channel
}

func (c *Exchange) unsubscribe(s *Subscription) {
	c.Lock()
	for i, subscriber := range c.subscriptors {
		if subscriber == s {
			n := i + 1
			left := c.subscriptors[:i]
			right := c.subscriptors[n:]
			c.subscriptors = append(left, right...)
		}
	}
	c.Unlock()
}

func (c *Exchange) subscribe(s *Subscription) {
	c.Lock()
	c.subscriptors = append(c.subscriptors, s)
	c.Unlock()
}

// Subscribe to channel
//
//     subscription := channel.Subscribe()
//     msg := <- sbuscription.C
//
func (c *Exchange) Subscribe() *Subscription {
	subscription := &Subscription{make(chan interface{}), c}
	c.subscribe(subscription)
	return subscription
}

// Subscriptors return the number of subscriptors
func (c *Exchange) Subscriptors() int {
	return len(c.subscriptors)
}

// Publish a message into the channel (Broadcast)
// It will go to all the subscriptions and send individually the message
// It will block until all the subscriptors recive the message
// You may want to launch this in its independent gorutine
//
//     go channel.Publish("msg")
//
func (c *Exchange) Publish(data interface{}) {
	c.Lock()
	for _, subscriber := range c.subscriptors {
		subscriber.C <- data
	}
	c.Unlock()
}

// Unbscribe the subscription
// This step is necessary to skip memory leaks
func (s *Subscription) Unsubscribe() {
	s.channel.unsubscribe(s)
}
