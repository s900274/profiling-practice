package main

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func metricHandler(w http.ResponseWriter, read *http.Request) {
	scale := time.Nanosecond
	du := float64(scale)
	duSuffix := scale.String()[1:]
	l := log.New(w, "metrics: ", log.Lmicroseconds)
	r := metrics.DefaultRegistry
	r.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			l.Printf("counter %s\n", name)
			l.Printf("  count:       %9d\n", metric.Count())
		case metrics.Gauge:
			l.Printf("gauge %s\n", name)
			l.Printf("  value:       %9d\n", metric.Value())
		case metrics.GaugeFloat64:
			l.Printf("gauge %s\n", name)
			l.Printf("  value:       %f\n", metric.Value())
		case metrics.Healthcheck:
			metric.Check()
			l.Printf("healthcheck %s\n", name)
			l.Printf("  error:       %v\n", metric.Error())
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			l.Printf("histogram %s\n", name)
			l.Printf("  count:       %9d\n", h.Count())
			l.Printf("  min:         %9d\n", h.Min())
			l.Printf("  max:         %9d\n", h.Max())
			l.Printf("  mean:        %12.2f\n", h.Mean())
			l.Printf("  stddev:      %12.2f\n", h.StdDev())
			l.Printf("  median:      %12.2f\n", ps[0])
			l.Printf("  75%%:         %12.2f\n", ps[1])
			l.Printf("  95%%:         %12.2f\n", ps[2])
			l.Printf("  99%%:         %12.2f\n", ps[3])
			l.Printf("  99.9%%:       %12.2f\n", ps[4])
		case metrics.Meter:
			m := metric.Snapshot()
			l.Printf("meter %s\n", name)
			l.Printf("  count:       %9d\n", m.Count())
			l.Printf("  1-min rate:  %12.2f\n", m.Rate1())
			l.Printf("  5-min rate:  %12.2f\n", m.Rate5())
			l.Printf("  15-min rate: %12.2f\n", m.Rate15())
			l.Printf("  mean rate:   %12.2f\n", m.RateMean())
		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			l.Printf("timer %s\n", name)
			l.Printf("  count:       %9d\n", t.Count())
			l.Printf("  min:         %12.2f%s\n", float64(t.Min())/du, duSuffix)
			l.Printf("  max:         %12.2f%s\n", float64(t.Max())/du, duSuffix)
			l.Printf("  mean:        %12.2f%s\n", t.Mean()/du, duSuffix)
			l.Printf("  stddev:      %12.2f%s\n", t.StdDev()/du, duSuffix)
			l.Printf("  median:      %12.2f%s\n", ps[0]/du, duSuffix)
			l.Printf("  75%%:         %12.2f%s\n", ps[1]/du, duSuffix)
			l.Printf("  95%%:         %12.2f%s\n", ps[2]/du, duSuffix)
			l.Printf("  99%%:         %12.2f%s\n", ps[3]/du, duSuffix)
			l.Printf("  99.9%%:       %12.2f%s\n", ps[4]/du, duSuffix)
			l.Printf("  1-min rate:  %12.2f\n", t.Rate1())
			l.Printf("  5-min rate:  %12.2f\n", t.Rate5())
			l.Printf("  15-min rate: %12.2f\n", t.Rate15())
			l.Printf("  mean rate:   %12.2f\n", t.RateMean())
		}
		l.Printf("----------\n")
	})

	t := metrics.GetOrRegisterCounter("foo", nil)
	t.Inc(1)

	fmt.Println("YOYO")
}

func main() {
	c := metrics.NewCounter()
	metrics.Register("foo", c)
	c.Inc(47)

	g := metrics.NewGauge()
	metrics.Register("bar", g)
	g.Update(47)

	s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	metrics.Register("baz", h)
	h.Update(47)

	m := metrics.NewMeter()
	metrics.Register("quux", m)
	m.Mark(47)

	t := metrics.NewTimer()
	metrics.Register("bang", t)
	t.Time(func() {})
	t.Update(47)

	http.HandleFunc("/metrics", metricHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
