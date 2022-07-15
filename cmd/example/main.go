package main

import (
	"context"
	"fmt"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	venus_metrics "github.com/diwufeiwen/venus-metrics"
)

func recordWines() {
	ctx := context.TODO()
	for {
		ctx, _ = tag.New(
			ctx,
			tag.Upsert(venus_metrics.MinerID, string("f0128788")),
		)
		stats.Record(ctx, venus_metrics.NumberOfIsRoundWinner.M(1))
		time.Sleep(time.Second * 3)
	}
}

func recordComputeProofDuration() {
	ctx := context.TODO()
	for {
		done := venus_metrics.TimerSeconds(ctx, venus_metrics.ComputeProofDuration, "f0128788")
		time.Sleep(time.Second * 8)
		done()
	}
}

func main() {
	if err := venus_metrics.RegisterExporter(); err != nil {
		fmt.Printf("failed to register the exporter: %v", err)
		return
	}

	fmt.Println("start")

	go recordComputeProofDuration()

	go recordWines()

	fmt.Println("end")
}
