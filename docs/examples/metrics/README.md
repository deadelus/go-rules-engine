# Metrics Example

This example demonstrates how to monitor the rules engine's performance using the `MetricsCollector` interface.

## How it works

The engine allows you to inject a metrics collector that satisfies the `MetricsCollector` interface defined in `src/types.go`.

In this example, we implement a `FakePrometheusCollector` that simulates how you would integrate with a real monitoring system like **Prometheus** or **OpenTelemetry**.

### Metrics Tracked:
- **Rule Evaluations**: Records the duration and result of every single rule.
- **Engine Runs**: Records the total execution time and rule count.
- **Event Executions**: Records the duration and result of event handlers.

## Running the example

```bash
go run docs/examples/metrics/main.go
```

## Integrating with real Prometheus

To use real Prometheus, your implementation would look like this:

```go
type PrometheusCollector struct {
    ruleDuration *prometheus.HistogramVec
    runDuration  *prometheus.Histogram
}

func (c *PrometheusCollector) ObserveRuleEvaluation(name string, result bool, d time.Duration) {
    c.ruleDuration.WithLabelValues(name, strconv.FormatBool(result)).Observe(d.Seconds())
}

// ... and so on
```
