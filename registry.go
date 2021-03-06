package lft

import (
	"reflect"
	"sync"

	"github.com/rcrowley/go-metrics"
)

var DefaultRegistry = &Registry{}

// Registry implements Registry interface from go-metrics but
// uses sync.Map for cache
type Registry struct {
	cache sync.Map
}

// PullFrom adds all missing metrics from given registry and returns number
// of pulled metrics.
// Run periodically for compatibility during migrations.
func (r *Registry) PullFrom(source metrics.Registry) (n int) {
	source.Each(func(s string, i interface{}) {
		if _, loaded := r.cache.LoadOrStore(s, i); !loaded {
			n++
		}
	})
	return n
}

func (r *Registry) Each(fn func(string, interface{})) {
	r.cache.Range(func(key, value interface{}) bool {
		fn(key.(string), value)
		return true
	})
}

func (r *Registry) Get(name string) (m interface{}) {
	m, _ = r.cache.Load(name)
	return m
}

func (r *Registry) GetAll() map[string]map[string]interface{} {
	data := make(map[string]map[string]interface{})
	r.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch metric := i.(type) {
		case metrics.Counter:
			values["count"] = metric.Count()
		case metrics.Gauge:
			values["value"] = metric.Value()
		case metrics.GaugeFloat64:
			values["value"] = metric.Value()
		case metrics.Healthcheck:
			values["error"] = nil
			metric.Check()
			if err := metric.Error(); nil != err {
				values["error"] = metric.Error().Error()
			}
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = h.Count()
			values["min"] = h.Min()
			values["max"] = h.Max()
			values["mean"] = h.Mean()
			values["stddev"] = h.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
		case metrics.Meter:
			m := metric.Snapshot()
			values["count"] = m.Count()
			values["1m.rate"] = m.Rate1()
			values["5m.rate"] = m.Rate5()
			values["15m.rate"] = m.Rate15()
			values["mean.rate"] = m.RateMean()
		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = t.Count()
			values["min"] = t.Min()
			values["max"] = t.Max()
			values["mean"] = t.Mean()
			values["stddev"] = t.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
			values["1m.rate"] = t.Rate1()
			values["5m.rate"] = t.Rate5()
			values["15m.rate"] = t.Rate15()
			values["mean.rate"] = t.RateMean()
		}
		data[name] = values
	})
	return data
}

func (r *Registry) GetOrRegister(name string, m interface{}) (m1 interface{}) {
	m1, ok := r.cache.Load(name)
	if ok {
		return m1
	}
	if v := reflect.ValueOf(m); v.Kind() == reflect.Func {
		m = v.Call(nil)[0].Interface()
	}
	switch m.(type) {
	case metrics.Counter, metrics.Gauge, metrics.GaugeFloat64,
		metrics.Healthcheck, metrics.Histogram, metrics.Meter, metrics.Timer:
		m1, _ = r.cache.LoadOrStore(name, m)
	}
	return m1
}

func (r *Registry) Register(name string, m interface{}) error {
	if _, exists := r.cache.LoadOrStore(name, m); exists {
		return metrics.DuplicateMetric(name)
	}
	return nil
}

func (r *Registry) RunHealthchecks() {
	r.Each(func(s string, i interface{}) {
		if h, ok := i.(metrics.Healthcheck); ok {
			h.Check()
		}
	})
}

func (r *Registry) Unregister(name string) {
	r.cache.Delete(name)
}

func (r *Registry) UnregisterAll() {
	r.cache.Range(func(key, value interface{}) bool {
		r.cache.Delete(key)
		return true
	})
}
