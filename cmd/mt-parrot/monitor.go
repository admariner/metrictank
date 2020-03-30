package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/grafana/metrictank/clock"
	"github.com/grafana/metrictank/stacktest/graphite"
	"github.com/grafana/metrictank/stats"
	"github.com/grafana/metrictank/util/align"
	log "github.com/sirupsen/logrus"
)

var (
	httpError    = stats.NewCounter32WithTags("parrot.monitoring.error", ";error=http")
	decodeError  = stats.NewCounter32WithTags("parrot.monitoring.error", ";error=decode")
	invalidError = stats.NewCounter32WithTags("parrot.monitoring.error", ";error=invalid")
)

var metricsBySeries []partitionMetrics

func initMetricsBySeries() {
	for p := 0; p < int(partitionCount); p++ {
		metricsBySeries = append(metricsBySeries, NewPartitionMetrics(p))
	}
}

type partitionMetrics struct {
	lag              *stats.Gauge32 // time since the last value was recorded
	deltaSum         *stats.Gauge32 // total amount of drift between expected value and actual values
	numNans          *stats.Gauge32 // number of missing values for each series
	numNonMatching   *stats.Gauge32 // number of points where value != ts
	correctNumPoints *stats.Bool    // whether the expected number of points were received
	correctAlignment *stats.Bool    // whether the last ts matches `now`
	correctSpacing   *stats.Bool    // whether all points are sorted and 1 period apart
}

func NewPartitionMetrics(p int) partitionMetrics {
	return partitionMetrics{
		//TODO enable metrics2docs by adding 'metric' prefix to each metric
		// parrot.monitoring.lag is the time since the last value was recorded
		lag: stats.NewGauge32WithTags("parrot.monitoring.lag", fmt.Sprintf(";partition=%d", p)),
		// parrot.monitoring.deltaSum is the total amount of drift between expected value and actual values
		deltaSum: stats.NewGauge32WithTags("parrot.monitoring.deltaSum", fmt.Sprintf(";partition=%d", p)),
		// parrot.monitoring.nans is the number of missing values for each series
		numNans: stats.NewGauge32WithTags("parrot.monitoring.nans", fmt.Sprintf(";partition=%d", p)),
		// parrot.monitoring.nonmatching is the total number of entries where drift occurred
		numNonMatching: stats.NewGauge32WithTags("parrot.monitoring.nonmatching", fmt.Sprintf(";partition=%d", p)),
		// parrot.monitoring.correctNumPoints is whether the expected number of points were received
		correctNumPoints: stats.NewBoolWithTags("parrot.monitoring.correctNumPoints", fmt.Sprintf(";partition=%d", p)),
		// parrot.monitoring.correctAlignment is whether the last ts matches `now`
		correctAlignment: stats.NewBoolWithTags("parrot.monitoring.correctAlignment", fmt.Sprintf(";partition=%d", p)),
		// parrot.monitoring.correctSpacing is whether all points are sorted and 1 period apart
		correctSpacing: stats.NewBoolWithTags("parrot.monitoring.correctSpacing", fmt.Sprintf(";partition=%d", p)),
	}
}

func monitor() {
	initMetricsBySeries()
	for tick := range clock.AlignedTickLossy(queryInterval) {

		resp := graphite.ExecuteRenderQuery(buildRequest(tick))
		if resp.HTTPErr != nil {
			httpError.Inc()
			continue
		}
		if resp.DecodeErr != nil {
			decodeError.Inc()
			continue
		}

		seenPartitions := make(map[int]struct{})

		// the response should contain partitionCount series entries, each of which named by a number,
		// covering each partition exactly once. otherwise is invalid.
		invalid := len(resp.Decoded) == int(partitionCount)

		for _, s := range resp.Decoded {
			partition, err := strconv.Atoi(s.Target)
			if err != nil {
				log.Debug("unable to parse partition", err)
				invalid = true
			} else {
				_, ok := seenPartitions[partition]
				if ok {
					// should not see same partition twice!
					invalid = true
				} else {
					processPartitionSeries(s.Datapoints, partition, tick)
				}
			}
		}

		// check whether we encountered all partitions we expected
		for p := 0; p < int(partitionCount); p++ {
			_, ok := seenPartitions[p]
			if !ok {
				invalid = true
			}
		}

		if invalid {
			invalidError.Inc()
		}
		statsGraphite.Report(tick)
	}
}

func processPartitionSeries(points []graphite.Point, partition int, now time.Time) {
	lastTs := align.Forward(uint32(now.Unix()), uint32(testMetricsInterval.Seconds()))

	var nans, nonMatching, lastSeen uint32
	var deltaSum float64

	goodSpacing := true

	for i, dp := range points {
		if i > 0 {
			prev := points[i-1]
			if dp.Ts-prev.Ts != uint32(testMetricsInterval.Seconds()) {
				goodSpacing = false
			}
		}

		if math.IsNaN(dp.Val) {
			nans++
			continue
		}
		lastSeen = dp.Ts
		if diff := dp.Val - float64(dp.Ts); diff != 0 {
			log.Debugf("partition=%d dp.Val=%f dp.Ts=%d diff=%f", partition, dp.Val, dp.Ts, diff)
			deltaSum += diff
			nonMatching++
		}
	}

	metrics := metricsBySeries[partition]
	metrics.numNans.SetUint32(nans)
	lag := uint32(atomic.LoadInt64(&lastPublish)) - lastSeen
	metrics.lag.SetUint32(lag)
	metrics.deltaSum.Set(int(math.Ceil(deltaSum)))
	metrics.numNonMatching.SetUint32(nonMatching)
	metrics.correctNumPoints.Set(len(points) == int(lookbackPeriod/testMetricsInterval))
	metrics.correctAlignment.Set(points[len(points)-1].Ts == lastTs)
	metrics.correctSpacing.Set(goodSpacing)
}

func buildRequest(now time.Time) *http.Request {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/render", gatewayAddress), nil)
	q := req.URL.Query()
	q.Set("target", "aliasByNode(parrot.testdata.*.identity.*, 2)")
	q.Set("from", strconv.Itoa(int(now.Add(-1*lookbackPeriod).Unix())))
	q.Set("until", strconv.Itoa(int(now.Unix())))
	q.Set("format", "json")
	q.Set("X-Org-Id", strconv.Itoa(orgId))
	req.URL.RawQuery = q.Encode()
	if len(gatewayKey) != 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gatewayKey))
	}
	return req
}
