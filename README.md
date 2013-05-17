    `import "."`

    Overview
    Index
    Examples

## Overview ▹

## Overview ▾

## Index ▹

## Index ▾

    func NewPubSub() *exchange
    type Subscription
    &nbsp_place_holder; &nbsp_place_holder; func (s *Subscription) Unsubscribe()

#### Examples

    NewPubSub

#### Package files

[pub_sub.go](/target/pub_sub.go)

## func [NewPubSub](/target/pub_sub.go?s=347:373#L14)

    
    func NewPubSub() *exchange

Create a new Publish/Subscribe Channel

▹ Example

▾ Example

Code:

    
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
    

## type [Subscription](/target/pub_sub.go?s=228:303#L8)

    
    type Subscription struct {
        C chan (interface{})
        // contains filtered or unexported fields
    }

### func (*Subscription) [Unsubscribe](/target/pub_sub.go?s=1727:1763#L70)

    
    func (s *Subscription) Unsubscribe()

Unbscribe the subscription This step is necessary to skip memory leaks

