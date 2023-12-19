package sync_methods

import (
	"sync"
	"testing"
)

func BenchmarkChannel_run(b *testing.B) {
	var counter int
	channel := make(chan bool, 1) // Buffered channel with capacity of 1

	b.ResetTimer()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			channel <- true // Acquire the "lock"
			counter++
			<-channel // Release the "lock"
			wg.Done()
		}()
	}
	wg.Wait()
	close(channel)
}
