// Metrics output to Datadog.
package datadog

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/rcrowley/go-metrics"
)

func Datadog(r metrics.Registry, d time.Duration, addr string) {
	DatadogWithConfig(r, d, Config{
		Addr: addr,
	})
}

func DatadogWithConfig(r metrics.Registry, d time.Duration, config Config) {
	c, err := statsd.New(config.Addr)
	if err != nil {
		log.Println(err)
	}
	for {
		if err := sh(r, c, config); nil != err {
			log.Println(err)
		}
		time.Sleep(d)
	}
}

type Config struct {
	Addr    string
	AppName string
}

func sh(r metrics.Registry, client *statsd.Client, config Config) error {
	r.Each(func(metricName string, i interface{}) {
		dd := parseMetricName(metricName)
		name := dd.name
		tags := append(dd.tags, baseTags(config)...)

		switch metric := i.(type) {
		case metrics.Counter:
			client.Count(name, metric.Count(), tags, 1)
		case metrics.Gauge:
			client.Gauge(name, float64(metric.Value()), tags, 1)
		case metrics.GaugeFloat64:
			client.Gauge(name, float64(metric.Value()), tags, 1)
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			client.Count(name+".count", h.Count(), tags, 1)
			client.Gauge(name+".min", float64(h.Min()), tags, 1)
			client.Gauge(name+".max", float64(h.Max()), tags, 1)
			client.Gauge(name+".mean", float64(h.Mean()), tags, 1)
			client.Gauge(name+".stddev", float64(h.StdDev()), tags, 1)
			client.Gauge(name+".p50", float64(ps[0]), tags, 1)
			client.Gauge(name+".p75", float64(ps[1]), tags, 1)
			client.Gauge(name+".p95", float64(ps[2]), tags, 1)
			client.Gauge(name+".p99", float64(ps[3]), tags, 1)
			client.Gauge(name+".p999", float64(ps[4]), tags, 1)
		case metrics.Meter:
			m := metric.Snapshot()
			client.Count(name+".count", m.Count(), tags, 1)
			client.Gauge(name+".1MinuteRate", float64(m.Rate1()), tags, 1)
			client.Gauge(name+".5MinuteRate", float64(m.Rate5()), tags, 1)
			client.Gauge(name+".15MinuteRate", float64(m.Rate15()), tags, 1)
			client.Gauge(name+".mean", float64(m.RateMean()), tags, 1)
		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			client.Count(name+".count", t.Count(), tags, 1)
			client.Gauge(name+".min", float64(t.Min()), tags, 1)
			client.Gauge(name+".max", float64(t.Max()), tags, 1)
			client.Gauge(name+".mean", float64(t.Mean()), tags, 1)
			client.Gauge(name+".stddev", float64(t.StdDev()), tags, 1)
			client.Gauge(name+".p50", float64(ps[0]), tags, 1)
			client.Gauge(name+".p75", float64(ps[1]), tags, 1)
			client.Gauge(name+".p95", float64(ps[2]), tags, 1)
			client.Gauge(name+".p99", float64(ps[3]), tags, 1)
			client.Gauge(name+".p999", float64(ps[4]), tags, 1)
			client.Gauge(name+".1MinuteRate", float64(t.Rate1()), tags, 1)
			client.Gauge(name+".5MinuteRate", float64(t.Rate5()), tags, 1)
			client.Gauge(name+".15MinuteRate", float64(t.Rate15()), tags, 1)
			client.Gauge(name+".meanRate", float64(t.RateMean()), tags, 1)
		}
	})
	return nil
}

func baseTags(config Config) []string {
	var baseTags []string
	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))
	if env != "" {
		baseTags = append(baseTags, fmt.Sprintf("environment:%v", env))
	}
	if config.AppName != "" {
		baseTags = append(baseTags, fmt.Sprintf("app:%v", config.AppName))
	}

	return baseTags
}

type metric struct {
	name string
	tags []string
}

var reName = regexp.MustCompile(`([a-zA-Z0-9]{1,}(\.[a-zA-Z0-9]{1,}){0,})`)
var reTags = regexp.MustCompile(`\[([a-zA-Z0-9\-\_]+(\:[a-zA-Z0-9\-\_]+){0,},{0,1})+\]`)

func parseMetricName(name string) metric {
	n := reName.FindString(name)
	tagStr := reTags.FindString(name)
	t := parseTags(tagStr)
	return metric{
		name: n,
		tags: t,
	}
}

func parseTags(tagStr string) []string {
	s1 := strings.TrimLeft(tagStr, "[")
	s2 := strings.TrimRight(s1, "]")

	if s2 == "" {
		return []string{}
	}

	strTags := strings.Split(s2, ",")
	return strTags
}
