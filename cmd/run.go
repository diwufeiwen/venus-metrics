package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/log/level"
	"github.com/urfave/cli/v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/graphite_exporter/collector"
	"github.com/prometheus/statsd_exporter/pkg/mapper"
	"github.com/prometheus/statsd_exporter/pkg/mappercache/lru"
	"github.com/prometheus/statsd_exporter/pkg/mappercache/randomreplacement"

	venus_metrics "github.com/diwufeiwen/venus-metrics"
)

// TODO(mr): this is copied verbatim from statsd_exporter/main.go. It should be a
// convenience function in mappercache, but that caused an import cycle.
func getCache(cacheSize int, cacheType string, registerer prometheus.Registerer) (mapper.MetricMapperCache, error) {
	var cache mapper.MetricMapperCache
	var err error
	if cacheSize == 0 {
		return nil, nil
	} else {
		switch cacheType {
		case "lru":
			cache, err = lru.NewMetricMapperLRUCache(registerer, cacheSize)
		case "random":
			cache, err = randomreplacement.NewMetricMapperRRCache(registerer, cacheSize)
		default:
			err = fmt.Errorf("unsupported cache type %q", cacheType)
		}

		if err != nil {
			return nil, err
		}
	}

	return cache, nil
}

func runCollector() error {
	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)

	c := collector.NewGraphiteCollector(logger, false, time.Minute*5)
	prometheus.MustRegister(c)

	metricMapper := &mapper.MetricMapper{}

	cacheSize := 1000
	cacheType := "lru"
	cache, err := getCache(cacheSize, cacheType, prometheus.DefaultRegisterer)
	if err != nil {
		os.Exit(1)
	}
	metricMapper.UseCache(cache)

	c.SetMapper(metricMapper)

	graphiteAddress := ":4568"
	tcpSock, err := net.Listen("tcp", graphiteAddress)
	if err != nil {
		level.Error(logger).Log("msg", "Error binding to TCP socket", "err", err)
		os.Exit(1)
	}
	go func() {
		for {
			conn, err := tcpSock.Accept()
			if err != nil {
				level.Error(logger).Log("msg", "Error accepting TCP connection", "err", err)
				continue
			}
			go func() {
				defer conn.Close()
				c.ProcessReader(conn)
			}()
		}
	}()

	udpAddress, err := net.ResolveUDPAddr("udp", graphiteAddress)
	if err != nil {
		level.Error(logger).Log("msg", "Error resolving UDP address", "err", err)
		os.Exit(1)
	}
	udpSock, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		level.Error(logger).Log("msg", "Error listening to UDP address", "err", err)
		os.Exit(1)
	}
	go func() {
		defer udpSock.Close()
		for {
			buf := make([]byte, 65536)
			chars, srcAddress, err := udpSock.ReadFromUDP(buf)
			if err != nil {
				level.Error(logger).Log("msg", "Error reading UDP packet", "from", srcAddress, "err", err)
				continue
			}
			go c.ProcessReader(bytes.NewReader(buf[0:chars]))
		}
	}()

	return nil
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "start venus-metrics daemon",
	Action: func(cctx *cli.Context) error {
		// Register all metric views
		if err := venus_metrics.RegisterView(); err != nil {
			return fmt.Errorf("failed to register the views: %v", err)
		}

		// run the Prometheus exporter as a scrape endpoint.
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			w.Write([]byte(`<html>
      <head><title>Graphite Exporter</title></head>
      <body>
      <h1>Graphite Exporter</h1>
      <p>Accepting plaintext Graphite samples over TCP and UDP on ` + ":4568" + `</p>
      <p><a href="` + "metrics" + `">Metrics</a></p>
      </body>
      </html>`))
		})
		server := &http.Server{
			Addr:    cctx.String("listen"),
			Handler: mux,
		}

		sigCh := make(chan os.Signal, 2)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			select {
			case sig := <-sigCh:
				log.Warnw("received shutdown", "signal", sig)
			case <-cctx.Done():
				log.Warn("received shutdown")
			}

			log.Info("Shutting down...")
			if err := server.Shutdown(context.TODO()); err != nil {
				log.Errorf("shutting down RPC server failed: %s", err)
			}
		}()

		// create collector
		err := runCollector()
		if err != nil {
			return fmt.Errorf("failed to create collector: %v", err)
		}

		// start the service
		log.Infof("start to listen %s", server.Addr)
		promlogConfig := &promlog.Config{}
		logger := promlog.New(promlogConfig)
		if err := web.ListenAndServe(server, "", logger); err != nil {
			level.Error(logger).Log("err", err)
			os.Exit(1)
		}

		log.Warn("Graceful shutdown successful")
		return nil
	},
}
