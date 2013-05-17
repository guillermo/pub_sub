package pub_sub

import (
	"fmt"
	"runtime"
)

func ExampleNewPubSub() {

	channel := NewPubSub()

	for i := 0; i < 2; i++ {
		goProc := i
		go func() {
			subscription := channel.Subscribe()
			for i := 0; i < 2; i++ {
				fmt.Println(goProc, <-subscription.C)
			}
			subscription.Unsubscribe()
		}()
	}
	runtime.Gosched()

	channel.Publish("hola")
	channel.Publish("adios")
	runtime.Gosched()
	channel.Stop()

	// Output:
	// 0 hola
	// 0 adios
	// 1 hola
	// 1 adios
}

func (c *Exchange) Subscribers() int {
	return len(c.subscribers)
}
