package go_pgx_prometheus_collector

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

type pgxStatsCollector struct {
	conn *pgxpool.Pool
	pms  []poolMetric
}

func (p *pgxStatsCollector) Describe(descs chan<- *prometheus.Desc) {
	for _, pm := range p.pms {
		descs <- pm.desc
	}
}

func (p *pgxStatsCollector) Collect(metrics chan<- prometheus.Metric) {
	stats := p.conn.Stat()
	for _, pm := range p.pms {
		metrics <- prometheus.MustNewConstMetric(pm.desc, pm.tp, pm.fn(stats))
	}
}

type poolMetric struct {
	desc *prometheus.Desc
	fn   func(stats *pgxpool.Stat) float64
	tp   prometheus.ValueType
}

// NewPgxCollector returns new collector for pgx pool with namespace `ns`.
// `ns` may be empty.
func NewPgxCollector(ns string,
	pgConn *pgxpool.Pool) prometheus.Collector {

	return &pgxStatsCollector{
		conn: pgConn,
		pms: []poolMetric{
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "pgx", "acquire_count"),
					"The cumulative count of successful acquires from "+
						"the pool.",
					nil, nil),
				fn: func(stats *pgxpool.Stat) float64 {
					return float64(stats.AcquireCount())
				},
				tp: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "pgx",
						"acquire_duration"),
					"The total duration of all successful acquires from "+
						"the pool.",
					nil, nil),
				fn: func(stats *pgxpool.Stat) float64 {
					return stats.AcquireDuration().Seconds()
				},
				tp: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "pgx", "acquire_conns"),
					"The number of currently acquired connections in the pool.",
					nil, nil),
				fn: func(stats *pgxpool.Stat) float64 {
					return float64(stats.AcquiredConns())
				},
				tp: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "pgx",
						"canceled_acquire_count"),
					"The cumulative count of acquires from the pool that "+
						"were canceled by a context.",
					nil, nil),
				fn: func(stats *pgxpool.Stat) float64 {
					return float64(stats.CanceledAcquireCount())
				},
				tp: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "pgx",
						"constructing_conns"),
					"The number of conns with construction in progress in "+
						"the pool.",
					nil, nil),
				fn: func(stats *pgxpool.Stat) float64 {
					return float64(stats.ConstructingConns())
				},
				tp: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "pgx",
						"empty_acquire_count"),
					"The cumulative count of successful acquires from the "+
						"pool that waited for a resource to be released or "+
						"constructed because the pool was empty.",
					nil, nil),
				fn: func(stats *pgxpool.Stat) float64 {
					return float64(stats.EmptyAcquireCount())
				},
				tp: prometheus.CounterValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "pgx", "idle_conns"),
					"The number of currently idle conns in the pool.",
					nil, nil),
				fn: func(stats *pgxpool.Stat) float64 {
					return float64(stats.IdleConns())
				},
				tp: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "pgx", "max_conns"),
					"The maximum size of the pool.",
					nil, nil),
				fn: func(stats *pgxpool.Stat) float64 {
					return float64(stats.MaxConns())
				},
				tp: prometheus.GaugeValue,
			},
			{
				desc: prometheus.NewDesc(
					prometheus.BuildFQName(ns, "pgx", "total_conns"),
					"The total number of resources currently in the pool. "+
						"The value is the sum of constructing_conns, "+
						"acquired_conns, and idle_conns.",
					nil, nil),
				fn: func(stats *pgxpool.Stat) float64 {
					return float64(stats.TotalConns())
				},
				tp: prometheus.GaugeValue,
			},
		},
	}
}
