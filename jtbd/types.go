// Package jtbd defines common types used across the JTBD testing framework.
package jtbd

import (
	"time"
)

// TestResults aggregates test execution results and metrics.
type TestResults struct {
	Results  []*ExecutionResult `json:"results"`
	Metrics  TestMetrics        `json:"metrics"`
	Duration time.Duration      `json:"duration"`
}
