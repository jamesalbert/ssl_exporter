package prober

import (
	"context"
	"fmt"

	"github.com/bmatcuk/doublestar/v2"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/jamesalbert/ssl_exporter/config"
)

// ProbeFile collects certificate metrics from local files
func ProbeFile(ctx context.Context, logger log.Logger, target string, module config.Module, registry *prometheus.Registry) error {
	errCh := make(chan error, 1)

	go func() {
		files, err := doublestar.Glob(target)
		if err != nil {
			errCh <- err
			return
		}

		if len(files) == 0 {
			errCh <- fmt.Errorf("No files found")
		} else {
			errCh <- collectFileMetrics(logger, files, registry)
		}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("context timeout, ran out of time")
	case err := <-errCh:
		return err
	}
}
