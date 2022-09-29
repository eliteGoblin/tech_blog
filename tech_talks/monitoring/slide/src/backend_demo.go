// START OMIT P1
var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	sinSignal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_sin_signal_seconds",
		Help: "Sin singal demo",
	})
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}
// END OMIT P1

// START OMIT P2
func recordSin() {
	go func() {
		for {
			sinSignal.Set(10.0 * math.Sin(
				2*math.Pi*float64(time.Now().Unix()%60.0)/60))
			time.Sleep(time.Second)
		}
	}()
}

func main() {
	recordMetrics()
	recordSin()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
// END OMIT P2