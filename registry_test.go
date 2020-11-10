package lft_test

import (
	"sort"
	"testing"

	"github.com/rcrowley/go-metrics"
	lft "github.com/remerge/go-lock_free_timer"
	"github.com/stretchr/testify/assert"
)

func TestRegistry_PullFrom(t *testing.T) {
	r := &lft.Registry{}
	b := metrics.NewRegistry()

	metrics.GetOrRegisterCounter("a", r)
	metrics.GetOrRegisterCounter("b", b)

	n := r.PullFrom(b)
	assert.Equal(t, 1, n)

	assert.NotNil(t, r.Get("a"))
	assert.NotNil(t, r.Get("b"))
}

func TestRegistry_Register(t *testing.T) {
	r := &lft.Registry{}
	assert.NoError(t, r.Register("a", metrics.NewCounter()))
	assert.NoError(t, r.Register("b", metrics.NewCounter()))
	assert.Error(t, r.Register("a", metrics.NewCounter()))

	var names []string
	r.Each(func(s string, i interface{}) {
		switch i.(type) {
		case metrics.Counter:
		default:
			t.Errorf("%s:%v is not counter", s, i)
		}
		names = append(names, s)
	})
	sort.Strings(names)
	assert.Equal(t, []string{"a", "b"}, names)
}

func TestRegistry_Each(t *testing.T) {
	r := &lft.Registry{}

	metrics.GetOrRegisterCounter("a", r)
	metrics.GetOrRegisterCounter("b", r)
	metrics.GetOrRegisterCounter("a", r)

	var names []string
	r.Each(func(s string, i interface{}) {
		switch i.(type) {
		case metrics.Counter:
		default:
			t.Errorf("%s:%v is not counter", s, i)
		}
		names = append(names, s)
	})
	sort.Strings(names)
	assert.Equal(t, []string{"a", "b"}, names)
}

func BenchmarkRegistry_GetOrRegister(b *testing.B) {
	b.Run(`builtin`, func(b *testing.B) {
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				metrics.GetOrRegisterCounter("a", metrics.DefaultRegistry)
			}
		})
	})
	b.Run(`lft`, func(b *testing.B) {
		r := &lft.Registry{}
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				metrics.GetOrRegisterCounter("a", r)
			}
		})
	})
}
