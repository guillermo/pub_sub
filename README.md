# Overview

Publish Subscribe with go channels

This library provie the pub/sub pattern to the go routines

# Index

    type Exchange
    func NewPubSub() *Exchange
    func (e *Exchange) Publish(data interface{})
    func (e *Exchange) Subscribe() *Subscription
    type Subscription
    func (s *Subscription) Unsubscribe()


# Package files

[pub_sub.go](pub_sub.go)

## type **Exchange**


    type Exchange struct {
        sync.Mutex
        // contains filtered or unexported fields
    }

### func **NewPubSub**


####    func NewPubSub() *Exchange

Create a new Publish/Subscribe Channel

Example:


    // We create a exchange channel
    channel := NewPubSub()

    // If we publish, the message will be lost
    channel.Publish("Hello World")

    // So lets do a few subscribers
    for i := 0; i < 3; i++ {

        go func() {
            // We need a subscription
            subscription := channel.Subscribe()
            for i := 0; i < 2; i++ {
                // And lets wait for messages
                msg := <-subscription.C
                // Whilte the messages are being to all the subcriptors, the publisher is block
                // Is responsability of the reciver to don't take too much time
                go func(msg interface{}) {
                    fmt.Println(msg)
                }(msg)
            }
            subscription.Unsubscribe()
        }()
    }

    runtime.Gosched() //Need to allocate the subscribers
    // And start publish messages
    for channel.Subscribers() != 0 {
        _ = <-time.After(time.Second)
        channel.Publish("hey")
    }

    Outputo

    hey
    hey
    hey
    hey
    hey
    hey


### func (*Exchange) **Publish**


    func (e *Exchange) Publish(data interface{})

Publish a message into the channel (Broadcast) It will go to all the
subscriptions and send individually the message It will block until all the
subscriptors recive the message You may want to launch this in its independent
gorutine


    go channel.Publish("msg")


### func (*Exchange) **Subscribe**


    func (e *Exchange) Subscribe() *Subscription

Subscribe to channel


    subscription := channel.Subscribe()
    msg := <- sbuscription.C


### func (*Exchange) **Subscriptors**


    func (e *Exchange) Subscriptors() [int](/pkg/builtin/#int)

Subscriptors return the number of subscriptors

## type **Subscription**


    type Subscription struct {
        C chan (interface{})
        // contains filtered or unexported fields
    }

### func (*Subscription)


    func (s *Subscription) Unsubscribe()

Unbscribe the subscription This step is necessary to skip memory leaks

