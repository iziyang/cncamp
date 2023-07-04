package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	numProducers := 2 // 设置生产者数量
	numConsumers := 3 // 设置消费者数量
	var consumerChannels []chan int

	for i := 0; i < numConsumers; i++ {
		consumerChannels = append(consumerChannels, make(chan int, 10))
	}

	done := make(chan bool)
	var wg sync.WaitGroup

	// consumer
	consumer := func(id int, messages <-chan int) {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-done:
				fmt.Printf("consumer %d: interrupt...\n", id)
				return
			case msg := <-messages:
				fmt.Printf("consumer %d: received message: %d\n", id, msg)
			case <-ticker.C:
				fmt.Printf("consumer %d: no message\n", id)
			}
		}
	}

	// producer
	producer := func(id int) {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			// 把消息发送给所有消费者的 channel
			for _, ch := range consumerChannels {
				ch <- i
			}
			time.Sleep(500 * time.Millisecond)
		}
	}

	// 启动消费者 Goroutines
	wg.Add(numConsumers)
	for i := 0; i < numConsumers; i++ {
		go consumer(i, consumerChannels[i])
	}

	// 启动生产者 Goroutines
	wg.Add(numProducers)
	for i := 0; i < numProducers; i++ {
		go producer(i)
	}

	wg.Wait()
	close(done)
	fmt.Println("main process exit!")
}
