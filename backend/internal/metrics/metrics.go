// Package metrics defines the app's custom Prometheus metrics (as opposed to
// the generic HTTP metrics echoprometheus already provides for the API
// server). All of it is exposed on the admin port's /metrics endpoint — see
// server.NewHealthServer.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"

	"todoapp/internal/models"
)

const namespace = "todoapp"

var (
	TodosCreatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "todos_created_total",
		Help:      "Total number of todos ever created.",
	})

	TodosOpen = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "todos_open",
		Help:      "Current number of todos that are not completed.",
	})

	TodosClosed = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "todos_closed",
		Help:      "Current number of todos that are completed.",
	})
)

func init() {
	prometheus.MustRegister(TodosCreatedTotal, TodosOpen, TodosClosed)
}

// SeedFromDB sets the open/closed gauges to the database's actual current
// counts. Gauges start at zero on every process restart, so without this a
// restart would under-report until enough toggles/creates/deletes happened
// to correct it — this makes the metrics accurate immediately on startup.
func SeedFromDB(db *gorm.DB) error {
	var openCount, closedCount int64

	if err := db.Model(&models.Todo{}).Where("completed = ?", false).Count(&openCount).Error; err != nil {
		return err
	}
	if err := db.Model(&models.Todo{}).Where("completed = ?", true).Count(&closedCount).Error; err != nil {
		return err
	}

	TodosOpen.Set(float64(openCount))
	TodosClosed.Set(float64(closedCount))
	return nil
}
