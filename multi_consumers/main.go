package main

import (
	"fmt"
	"time"
)

func main() {
	numConsumers := 3 // 设置消费者数量
	var consumerChannels []chan int

	for i := 0; i < numConsumers; i++ {
		consumerChannels = append(consumerChannels, make(chan int, 10))
	}

	done := make(chan bool)

	// consumer
	consumer := func(id int, messages <-chan int) {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-done:
				fmt.Printf("consumer %d: interrupt...\n", id)
				return
			case msg := <-messages:
				fmt.Printf("consumer %d: received message: %d\n", id, msg)
			default:
				fmt.Printf("consumer %d: no message\n", id)
			}
		}
	}

	// 启动消费者 Goroutines
	for i := 0; i < numConsumers; i++ {
		go consumer(i, consumerChannels[i])
	}

	// producer
	for i := 0; i < 10; i++ {
		// 把消息发送给所有消费者的 channel
		for _, ch := range consumerChannels {
			ch <- i
		}
	}
	time.Sleep(5001 * time.Millisecond)
	close(done)
	time.Sleep(1001 * time.Millisecond)
	fmt.Println("main process exit!")
}
