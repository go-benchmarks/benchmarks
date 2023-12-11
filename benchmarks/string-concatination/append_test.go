package string_concatination

import (
	"testing"
)

func BenchmarkSimpleAppend_write(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s string
		s = s + "a"
	}
}

func BenchmarkSimpleAppend_read(b *testing.B) {
	var s string
	for i := 0; i < readCount; i++ {
		s = s + "a"
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s
	}
}
