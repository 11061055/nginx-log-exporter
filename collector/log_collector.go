package metric

import (
	"fmt"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/hpcloud/tail"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/11061055/nginx-log-exporter/config"
)

// Collector is a struct containing pointers to all metrics that should be
// exposed to Prometheus
type Collector struct {
	countTotal      *prometheus.CounterVec
	bytesTotal      *prometheus.CounterVec
	upstreamSeconds *prometheus.HistogramVec
	responseSeconds *prometheus.HistogramVec

	staticValues    []string
	dynamicLabels   []string
	dynamicValueLen int

	cfg    *config.AppConfig
}

func NewCollector(cfg *config.AppConfig) *Collector {
	staticLabels, staticValues := cfg.StaticLabelValues()
	dynamicLabels := cfg.DynamicLabels()

	labels := append(staticLabels, dynamicLabels...)
    bucket := prometheus.LinearBuckets(cfg.HistogramBuckets.Start, cfg.HistogramBuckets.Step, cfg.HistogramBuckets.Num)

	return &Collector{
		countTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: cfg.Name,
			Name:      "http_response_count_total",
			Help:      "Amount of processed HTTP requests",
		}, labels),

		bytesTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: cfg.Name,
			Name:      "http_response_size_bytes",
			Help:      "Total amount of transferred bytes",
		}, labels),

		upstreamSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Name,
			Name:      "http_upstream_time_seconds",
			Help:      "Time needed by upstream servers to handle requests",
			Buckets:   bucket,
		}, labels),

		responseSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Name,
			Name:      "http_response_time_seconds",
			Help:      "Time needed by NGINX to handle requests",
			Buckets:   bucket,
		}, labels),

		staticValues:    staticValues,
		dynamicLabels:   dynamicLabels,
		dynamicValueLen: len(dynamicLabels),

		cfg:    cfg,
	}
}

func (c *Collector) Run() {
	c.cfg.Prepare()

	// register to prometheus
	prometheus.MustRegister(c.countTotal)
	prometheus.MustRegister(c.bytesTotal)
	prometheus.MustRegister(c.upstreamSeconds)
	prometheus.MustRegister(c.responseSeconds)

	for _, f := range c.cfg.SourceFiles {
		t, err := tail.TailFile(f, tail.Config{
		    //ReOpen: true,
			Follow: true,
			Poll:   true,
		})

		if err != nil {
			log.Panic(err)
		}

		go func() {
			for line := range t.Lines {

			    var mp map[string]string

            	err = json.Unmarshal([]byte(line.Text), &mp)
            	if err != nil {
					fmt.Printf("error while parsing line '%s': %s", line.Text, err)
            		continue
            	}

				dynamicValues := make([]string, c.dynamicValueLen)

				for i, label := range c.dynamicLabels {

					if s, ok := mp[label]; ok {
						dynamicValues[i] = c.formatValue(label, s)
					}
				}

				labelValues := append(c.staticValues, dynamicValues...)

				c.countTotal.WithLabelValues(labelValues...).Inc()

				if bytes, ok := mp["body_bytes_sent"]; ok {

				    if b, err := strconv.ParseFloat(bytes, 32); err == nil {

					    c.bytesTotal.WithLabelValues(labelValues...).Add(b)
					}
				}

				if upstreamTime, ok := mp["upstream_response_time"]; ok {

				    if u, err := strconv.ParseFloat(upstreamTime, 32); err == nil {

					    c.upstreamSeconds.WithLabelValues(labelValues...).Observe(u)
					}
				}

				if responseTime, ok := mp["request_time"]; ok {

				    if r, err := strconv.ParseFloat(responseTime, 32); err == nil {

					    c.responseSeconds.WithLabelValues(labelValues...).Observe(r)
					}
				}
			}
		}()
	}
}

func (c *Collector) formatValue(label, value string) string {
	replacements, ok := c.cfg.RelabelConfig.Replacements[label]
	if !ok {
		return value
	}

	for _, replacement := range replacements {

	    if replacement.Trims != nil {

	        for _, trim := range replacement.Trims {
	            arr := strings.Split(value, trim.Sep)
	            if len(arr) > trim.Idx {
	                value = arr[trim.Idx]
	            }
	        }
	    }

	    if replacement.Repace != nil {

	        for _, target := range replacement.Repace {
		        if target.Regexp().MatchString(value) {
		            value = target.Regexp().ReplaceAllString(value, target.Value)
		        }
	        }
	    }
	}

	return value
}
