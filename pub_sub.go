// License: BSD

/*
Package pub_sub provides Publish Subscribe patther to gorutines
It is an easy way to multiplex messages to many recivers/channels

To define a publish/subscribe channel you first need to create the exchange:

	exchange = NewPubSub()

Once you created the exchange the subscribers need to subscribe to that exchange:

	subscription = exchange.Subscribe()

The next is the subscribers waiting for messages:

	for {
	  msg := <- subscribers.C
	  fmt.Println(msg)
	}

Now you can publish events with:

	exchange.Publish("hello")

Publish send a message to every of his subscribers channels.

Once you want to stop you have to cancel your subscription:

	subscription.Unsubscribe()

And to stop the exchange:

	exchange.Stop()

Once the exchange is stop, you will need to create a new exchange.
*/
package pub_sub

// Exchange represent the main point to deliver messages. You get one with NewPubSub()
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
				go func() { subscriber.C <- msg }()
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

// Subscribe generates a new subscription to the exchange. Remember to call unsubscribe.
//
//     subscription := exchange.Subscribe()
//     msg := <- sbuscription.C
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
	close(s.C)
}

// Stop the gorutine for the Exchange
func (e *Exchange) Stop() {
	e.closeChannel <- true
}
