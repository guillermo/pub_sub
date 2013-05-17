package pub_sub

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {

	channel := NewPubSub()
	subscribeChannel := channel.Subscribe()

	go channel.Publish("hola")
	data := <-subscribeChannel.C
	if data != "hola" {
		t.Error("Didn't work")
	}
	subscribeChannel.Unsubscribe()

}

func ExampleNewPubSub() {

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

	// Output:
	// hey
	// hey
	// hey
	// hey
	// hey
	// hey

}
