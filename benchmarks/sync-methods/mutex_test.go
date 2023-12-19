package sync_methods

import (
	"sync"
	"testing"
)

func BenchmarkMutex_run(b *testing.B) {
	var counter int
	var mutex sync.Mutex

	b.ResetTimer()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			mutex.Lock()
			counter++
			mutex.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
}
