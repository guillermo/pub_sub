# Overview

Publish Subscribe with go channels

This library provie the pub/sub pattern to the go routines

# Index

    type Exchange
    &nbsp_place_holder; &nbsp_place_holder; func NewPubSub() *Exchange
    &nbsp_place_holder; &nbsp_place_holder; func (c *Exchange) Publish(data interface{})
    &nbsp_place_holder; &nbsp_place_holder; func (c *Exchange) Subscribe() *Subscription
    &nbsp_place_holder; &nbsp_place_holder; func (c *Exchange) Subscriptors() int
    type Subscription
    &nbsp_place_holder; &nbsp_place_holder; func (s *Subscription) Unsubscribe()


# Package files

[pub_sub.go](/target/pub_sub.go)

## type [Exchange](/target/pub_sub.go?s=157:225#L2)


    type Exchange struct {
        [sync](/pkg/sync/).[Mutex](/pkg/sync/#Mutex)
        // contains filtered or unexported fields
    }

### func [NewPubSub](/target/pub_sub.go?s=346:372#L13)


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
    for channel.Subscriptors() != 0 {
        _ = <-time.After(time.Second)
        channel.Publish("hey")
    }

    Output:

    hey
    hey
    hey
    hey
    hey
    hey


### func (*Exchange) [Publish](/target/pub_sub.go?s=1509:1553#L60)


    func (c *Exchange) Publish(data interface{})

Publish a message into the channel (Broadcast) It will go to all the
subscriptions and send individually the message It will block until all the
subscriptors recive the message You may want to launch this in its independent
gorutine


    go channel.Publish("msg")


### func (*Exchange) [Subscribe](/target/pub_sub.go?s=948:992#L42)


    func (c *Exchange) Subscribe() *Subscription

Subscribe to channel


    subscription := channel.Subscribe()
    msg := <- sbuscription.C


### func (*Exchange) [Subscriptors](/target/pub_sub.go?s=1154:1191#L49)


    func (c *Exchange) Subscriptors() [int](/pkg/builtin/#int)

Subscriptors return the number of subscriptors

## type [Subscription](/target/pub_sub.go?s=227:302#L7)


    type Subscription struct {
        C chan (interface{})
        // contains filtered or unexported fields
    }

### func (*Subscription) [Unsubscribe](/target/pub_sub.go?s=1729:1765#L70)


    func (s *Subscription) Unsubscribe()

Unbscribe the subscription This step is necessary to skip memory leaks

