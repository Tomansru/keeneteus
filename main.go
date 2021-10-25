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
	networkStat = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "keeneteus_network",
		Help: "Used traffic per interface",
	}, []string{"interface", "rxtx"})
	devicesStat = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "keeneteus_devices",
		Help: "Used traffic per devices",
	}, []string{"device", "rxtx"})
	devicesRssiStat = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "keeneteus_devices_rssi",
		Help: "Used traffic per devices",
	}, []string{"device"})
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

	prometheus.MustRegister(cpuLoad, memUsage, uptimeStat, networkStat, devicesStat, devicesRssiStat)

	go func() {
		var dev string
		var upt int
		var m keenetic_api.Metrics
		var i keenetic_api.InterfaceStat
		i.SetInterfaces([]keenetic_api.Eth{
			{Name: "DOM.RU", Code: "GigabitEthernet0/Vlan4"},
			{Name: "Mishek.NET", Code: "GigabitEthernet1"},
			{Name: "WGHetzner", Code: "Wireguard0"}})
		i.SetDevices([]keenetic_api.Eth{
			{Name: "StanislavPC", Code: "18:c0:4d:64:4c:1e"},
			{Name: "MacBook 2015", Code: "a8:66:7f:2e:4a:d2"},
			{Name: "OnePlus6T", Code: "c0:ee:fb:4c:60:fd"},
			{Name: "Multicast", Code: "multicast"},
			{Name: "Others", Code: "others"}})
		for range time.Tick(time.Second * 3) {
			if err = kApi.Metric(&i); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, v := range i.Show.Interface.Stat {
				networkStat.WithLabelValues(v.InterfaceName, "rx").Set(float64(v.Rxbytes))
				networkStat.WithLabelValues(v.InterfaceName, "tx").Set(float64(v.Txbytes))
			}

			for k, v := range i.Show.Ip.Hotspot.Chart.Bar {
				dev = i.GetDeviceName(k)
				for _, v2 := range v.Bars {
					if v2.Attribute == "" || len(v2.Data) == 0 {
						break
					}
					devicesStat.WithLabelValues(dev, v2.Attribute[:2]).Add(float64(v2.Data[len(v2.Data)-1].V))
				}
			}

			if err = kApi.Metric(&m); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, v := range m.Show.Ip.Hotspot.Host {
				devicesRssiStat.WithLabelValues(v.Name).Set(float64(v.Rssi))
				devicesRssiStat.WithLabelValues(v.Name).Set(float64(v.Rssi))
			}

			upt, _ = strconv.Atoi(m.Show.System.Uptime)
			uptimeStat.Set(float64(upt))
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
