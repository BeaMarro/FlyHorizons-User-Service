package metrics

import (
	"flyhorizons-userservice/internal/health"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetricsRoutes(router *gin.Engine, dbCheck health.DatabaseCheck, rabbitMQCheck health.RabbitMQCheck) {
	dbHealthGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mssql_db_health",
		Help: "Database health status: 1 for up, 0 for down",
	})

	rabbitMQHealthGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rabbitmq_health",
		Help: "RabbitMQ health status: 1 for up, 0 for down",
	})

	// Register both metrics
	prometheus.MustRegister(dbHealthGauge)
	prometheus.MustRegister(rabbitMQHealthGauge)

	go func() {
		for {
			if dbCheck.Pass() {
				dbHealthGauge.Set(1)
			} else {
				dbHealthGauge.Set(0)
			}

			if rabbitMQCheck.Pass() {
				rabbitMQHealthGauge.Set(1)
			} else {
				rabbitMQHealthGauge.Set(0)
			}

			time.Sleep(10 * time.Second) // adjust interval as needed
		}
	}()

	// Expose /metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
