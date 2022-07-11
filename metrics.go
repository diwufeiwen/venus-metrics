package venus_metrics

import (
	"context"
	"fmt"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// Distribution(区间聚合)
var defaultMillisecondsDistribution = view.Distribution(0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 3000, 4000, 5000, 7500, 10000, 20000, 50000, 100000)

// Measures
var (
	// venus-miner
	MinerGetBaseInfoDuration   = stats.Float64("miner/getbaseinfo_ms", "Duration of getStateNonce in mpool", stats.UnitMilliseconds)
	MinerComputeTicketDuration = stats.Float64("miner/computeticket_ms", "Duration of getStateNonce in mpool", stats.UnitMilliseconds)
	MinerIsRoundWinnerDuration = stats.Float64("miner/isroundwinner_ms", "Duration of getStateNonce in mpool", stats.UnitMilliseconds)
	MinerComputeProofDuration  = stats.Float64("miner/computeproof_ms", "Duration of getStateNonce in mpool", stats.UnitMilliseconds)

	// venus-messager

)

var (
	KeyMethod, _ = tag.NewKey("method")
	KeyStatus, _ = tag.NewKey("status")
	KeyError, _  = tag.NewKey("error")
)

var (
	MinerGetBaseInfoDurationView = &view.View{
		Measure:     MinerGetBaseInfoDuration,
		Aggregation: defaultMillisecondsDistribution,
	}
	MinerComputeTicketDurationView = &view.View{
		Measure:     MinerComputeTicketDuration,
		Aggregation: defaultMillisecondsDistribution,
	}
	MinerIsRoundWinnerDurationView = &view.View{
		Measure:     MinerIsRoundWinnerDuration,
		Aggregation: defaultMillisecondsDistribution,
	}
	MinerComputeProofDurationView = &view.View{
		Measure:     MinerComputeProofDuration,
		Aggregation: defaultMillisecondsDistribution,
	}
)

var MinerNodeViews = append([]*view.View{
	MinerGetBaseInfoDurationView,
	MinerComputeTicketDurationView,
	MinerIsRoundWinnerDurationView,
	MinerComputeProofDurationView,
})

// SinceInMilliseconds returns the duration of time since the provide time as a float64.
func SinceInMilliseconds(startTime time.Time) float64 {
	return float64(time.Since(startTime).Nanoseconds()) / 1e6
}

// Timer is a function stopwatch, calling it starts the timer,
// calling the returned function will record the duration.
func Timer(ctx context.Context, m *stats.Float64Measure) func() {
	start := time.Now()
	return func() {
		stats.Record(ctx, m.M(SinceInMilliseconds(start)))
	}
}

func RegisterView() error {
	err := view.Register(MinerNodeViews...)
	if err != nil {
		return fmt.Errorf("failed to the miner-node view: %s", err)
	}

	return nil
}
