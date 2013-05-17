// License: BSD

// Publish Subscribe with go exchanges
//
// This library provie the pub/sub pattern to the go routines
package pub_sub

import ()

type Exchange struct {
	subscribers     []*Subscription
	subscribeChan   chan *Subscription
	unSubscribeChan chan *Subscription
	publishChan     chan interface{}
	closeChannel    chan bool
}

func reciver(e *Exchange) {
	for {
		select {
		case s := <-e.subscribeChan:
			e.subscribers = append(e.subscribers, s)
		case s := <-e.unSubscribeChan:
			for i, subscriber := range e.subscribers {
				if subscriber == s {
					n := i + 1
					left := e.subscribers[:i]
					right := e.subscribers[n:]
					e.subscribers = append(left, right...)
				}
			}
		case msg := <-e.publishChan:
			for _, subscriber := range e.subscribers {
				subscriber.C <- msg
			}
		case _ = <-e.closeChannel:
			break
		}
	}
}

// NewPubSub returns a new exchange
func NewPubSub() *Exchange {
	exchange := new(Exchange)
	exchange.subscribeChan = make(chan *Subscription)
	exchange.unSubscribeChan = make(chan *Subscription)
	exchange.closeChannel = make(chan bool)
	exchange.publishChan = make(chan interface{}, 2000)

	go reciver(exchange)
	return exchange
}

type Subscription struct {
	C        chan interface{}
	exchange *Exchange
}

// Subscribe to exchange
//
//     subscription := exchange.Subscribe()
//     msg := <- sbuscription.C
//
func (e *Exchange) Subscribe() *Subscription {
	subscription := &Subscription{make(chan interface{}), e}
	e.subscribeChan <- subscription
	return subscription
}

// Publish a message into the exchange (Broadcast)
// It will go to all the subscriptions and send individually the message
// It will block until all the subscribers recive the message
// You may want to launch this in its independent gorutine
//
//     go exchange.Publish("msg")
//
func (e *Exchange) Publish(msg interface{}) {
	e.publishChan <- msg
}

// Unbscribe the subscription
// This step is necessary to skip memory leaks
func (s *Subscription) Unsubscribe() {
	s.exchange.unSubscribeChan <- s
}

// Stop the gorutine for the Exchange
func (e *Exchange) Stop() {
	e.closeChannel <- true
}
