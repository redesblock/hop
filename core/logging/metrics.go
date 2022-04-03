package logging

import (
	"github.com/prometheus/client_golang/prometheus"
	m "github.com/redesblock/hop/core/metrics"
	"github.com/sirupsen/logrus"
)

type metrics struct {
	// all metrics fields must be exported
	// to be able to return them by Metrics()
	// using reflection
	ErrorCount prometheus.Counter
	WarnCount  prometheus.Counter
	InfoCount  prometheus.Counter
	DebugCount prometheus.Counter
	TraceCount prometheus.Counter
}

func newMetrics() (m metrics) {
	return metrics{
		ErrorCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "log_error_count",
			Help: "Number ERROR log messages.",
		}),
		WarnCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "log_warn_count",
			Help: "Number WARN log messages.",
		}),
		InfoCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "log_info_count",
			Help: "Number INFO log messages.",
		}),
		DebugCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "log_debug_count",
			Help: "Number DEBUG log messages.",
		}),
		TraceCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "log_trace_count",
			Help: "Number TRACE log messages.",
		}),
	}
}

func (l *logger) Metrics() []prometheus.Collector {
	return m.PrometheusCollectorsFromFields(l.metrics)
}

func (m metrics) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

func (m metrics) Fire(e *logrus.Entry) error {
	switch e.Level {
	case logrus.ErrorLevel:
		m.ErrorCount.Inc()
	case logrus.WarnLevel:
		m.WarnCount.Inc()
	case logrus.InfoLevel:
		m.InfoCount.Inc()
	case logrus.DebugLevel:
		m.DebugCount.Inc()
	case logrus.TraceLevel:
		m.TraceCount.Inc()
	}
	return nil
}
