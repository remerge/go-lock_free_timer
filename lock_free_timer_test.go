package lft

import (
	"testing"
	"time"

	"github.com/rcrowley/go-metrics"
)

func BenchmarkUpdate(b *testing.B) {
	t := metrics.NewTimer()
	for n := 0; n < b.N; n++ {
		t.Update(time.Duration(123 * time.Millisecond))
	}
}

func BenchmarkLockFreeUpdate(b *testing.B) {
	t := NewLockFreeTimer()
	for n := 0; n < b.N; n++ {
		t.Update(time.Duration(123 * time.Millisecond))
	}
}
