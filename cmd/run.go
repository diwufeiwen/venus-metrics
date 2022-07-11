package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"

	venus_metrics "github.com/diwufeiwen/venus-metrics"
)

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "start venus-metrics daemon",
	Action: func(cctx *cli.Context) error {
		// Register all metric views
		if err := venus_metrics.RegisterView(); err != nil {
			log.Fatalf("Failed to register the views: %v", err)
		}

		// run the Prometheus exporter as a scrape endpoint.
		mux := http.NewServeMux()
		mux.Handle("/metrics", venus_metrics.Exporter())
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

		// start the service
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			return err
		}

		log.Warn("Graceful shutdown successful")
		return nil
	},
}
