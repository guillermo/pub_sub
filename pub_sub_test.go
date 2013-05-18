package pub_sub

import (
	"fmt"
	"runtime"
)

func ExampleNewPubSub() {

	channel := NewPubSub()
	subscription := channel.Subscribe()
	go func() {
		msg := <-subscription.C
		fmt.Println(msg)
		subscription.Unsubscribe()
	}()

	runtime.Gosched()
	channel.Publish("hello")
	runtime.Gosched()

	// Output:
	// hello
}
