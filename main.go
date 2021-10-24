package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Tomansru/keeneteus/keenetic_api"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cpuLoad = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "keeneteus_cpu_load",
		Help: "Current load of the CPU",
	})
	memUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "keeneteus_mem_usage",
		Help: "Current mem usage",
	}, []string{"type"})
	uptimeStat = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "keeneteus_uptime",
		Help: "Uptime metric",
	})
)

func main() {
	var kUrl, kUser, kPasswd = os.Getenv("KeeneticUrl"),
		os.Getenv("KeeneticUser"),
		os.Getenv("KeeneticPassword")

	var kApi = keenetic_api.NewApi(kUrl, kUser, kPasswd)

	var err error
	if err = kApi.Auth(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	prometheus.MustRegister(cpuLoad, memUsage, uptimeStat)

	go func() {
		for range time.Tick(time.Second * 2) {
			var i keenetic_api.InterfaceStat
			if err = kApi.Metric(&i); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}()

	go func() {
		for range time.Tick(time.Second * 4) {
			var m keenetic_api.Metrics
			if err = kApi.Metric(&m); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			var i, _ = strconv.Atoi(m.Show.System.Uptime)
			uptimeStat.Set(float64(i))
			cpuLoad.Set(float64(m.Show.System.Cpuload))
			memUsage.WithLabelValues("total").Set(float64(m.Show.System.Memtotal))
			memUsage.WithLabelValues("cache").Set(float64(m.Show.System.Memcache))
			memUsage.WithLabelValues("free").Set(float64(m.Show.System.Memfree))
			memUsage.WithLabelValues("buffers").Set(float64(m.Show.System.Membuffers))
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	if err = http.ListenAndServe("0.0.0.0:2112", nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
