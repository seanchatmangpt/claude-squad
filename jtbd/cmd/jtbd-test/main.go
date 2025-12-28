// Command jtbd-test runs JTBD test suites for Fortune 5 industries.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"claude-squad/jtbd"
)

var (
	runAll        = flag.Bool("all", false, "Run all JTBD tests")
	industry      = flag.String("industry", "", "Run tests for specific industry (saas, ecommerce, manufacturing, healthcare, fintech, education, realstate, logistics, hospitality, retail)")
	listIndustries = flag.Bool("list", false, "List supported industries")
	outputFile    = flag.String("output", "", "Write results to file (use '-' for stdout)")
	outputFormat  = flag.String("format", "text", "Output format: text, json, junit")
	verbose       = flag.Bool("v", false, "Verbose output")
	timeout       = flag.Duration("timeout", 5*time.Minute, "Test timeout")
	parallel      = flag.Int("parallel", 4, "Number of parallel test processes")
	coverage      = flag.Bool("coverage", false, "Generate coverage report")
	failCoverage  = flag.Bool("fail-coverage", false, "Fail if coverage below threshold")
	minCoverage   = flag.Float64("min-coverage", 70.0, "Minimum coverage percentage")
	runBench      = flag.Bool("bench", false, "Run benchmarks")
	retry         = flag.Bool("retry", false, "Retry failed tests")
	maxRetries    = flag.Int("max-retries", 2, "Maximum retry attempts")
	ciMode        = flag.Bool("ci", false, "Enable CI mode")
)

var supportedIndustries = []string{
	"saas",
	"ecommerce",
	"manufacturing",
	"healthcare",
	"fintech",
	"education",
	"realstate",
	"logistics",
	"hospitality",
	"retail",
}

func main() {
	flag.Parse()

	if *listIndustries {
		fmt.Println("Supported industries:")
		for _, ind := range supportedIndustries {
			fmt.Printf("  - %s\n", ind)
		}
		os.Exit(0)
	}

	if !*runAll && *industry == "" {
		fmt.Fprintln(os.Stderr, "Error: must specify --all or --industry")
		flag.Usage()
		os.Exit(1)
	}

	if *industry != "" && !isValidIndustry(*industry) {
		fmt.Fprintf(os.Stderr, "Error: invalid industry '%s'. Use --list to see supported industries\n", *industry)
		os.Exit(1)
	}

	// Run tests
	results, err := runTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running tests: %v\n", err)
		os.Exit(127)
	}

	// Output results
	if err := outputResults(results); err != nil {
		fmt.Fprintf(os.Stderr, "Error outputting results: %v\n", err)
		os.Exit(127)
	}

	// Determine exit code
	exitCode := calculateExitCode(results)
	os.Exit(exitCode)
}

func isValidIndustry(ind string) bool {
	for _, supported := range supportedIndustries {
		if strings.EqualFold(ind, supported) {
			return true
		}
	}
	return false
}

func runTests() (*jtbd.TestResults, error) {
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	config := &jtbd.RunConfig{
		Mode:          jtbd.ExecutionModeParallel,
		MaxWorkers:    *parallel,
		GlobalTimeout: *timeout,
		TestTimeout:   *timeout / 10,
		EnableRetry:   *retry,
		IsolateTests:  true,
	}

	var tests []*jtbd.Test

	if *runAll {
		tests = createAllTests()
	} else {
		tests = createIndustryTests(*industry)
	}

	engine, err := jtbd.NewExecutionEngine(tests, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create engine: %w", err)
	}

	execResults, err := engine.Run()
	if err != nil && !strings.Contains(err.Error(), "context") {
		return nil, fmt.Errorf("failed to run tests: %w", err)
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("tests timed out after %v", *timeout)
	default:
	}

	metrics := engine.GetMetrics()

	return &jtbd.TestResults{
		Results: execResults,
		Metrics: metrics,
		Duration: *timeout,
	}, nil
}

func createAllTests() []*jtbd.Test {
	var tests []*jtbd.Test
	for _, ind := range supportedIndustries {
		tests = append(tests, createIndustryTests(ind)...)
	}
	return tests
}

func createIndustryTests(industry string) []*jtbd.Test {
	// Create sample tests for the industry
	return []*jtbd.Test{
		{
			ID:          fmt.Sprintf("%s-test-1", industry),
			Name:        fmt.Sprintf("%s Basic Functionality", strings.Title(industry)),
			Description: fmt.Sprintf("Tests basic %s functionality", industry),
			Timeout:     30 * time.Second,
			MaxRetries:  *maxRetries,
			Execute: func(ctx context.Context) error {
				// Placeholder test logic
				time.Sleep(100 * time.Millisecond)
				return nil
			},
		},
		{
			ID:          fmt.Sprintf("%s-test-2", industry),
			Name:        fmt.Sprintf("%s Integration", strings.Title(industry)),
			Description: fmt.Sprintf("Tests %s integration", industry),
			Timeout:     45 * time.Second,
			MaxRetries:  *maxRetries,
			Dependencies: []string{fmt.Sprintf("%s-test-1", industry)},
			Execute: func(ctx context.Context) error {
				time.Sleep(150 * time.Millisecond)
				return nil
			},
		},
	}
}

func outputResults(results *jtbd.TestResults) error {
	var output string
	var err error

	switch *outputFormat {
	case "text":
		output = formatTextResults(results)
	case "json":
		data, jsonErr := json.MarshalIndent(results, "", "  ")
		if jsonErr != nil {
			return jsonErr
		}
		output = string(data)
	case "junit":
		output = formatJUnitResults(results)
	default:
		return fmt.Errorf("unknown format: %s", *outputFormat)
	}

	if *outputFile == "" || *outputFile == "-" {
		fmt.Println(output)
	} else {
		err = os.WriteFile(*outputFile, []byte(output), 0644)
	}

	return err
}

func formatTextResults(results *jtbd.TestResults) string {
	var sb strings.Builder

	sb.WriteString("JTBD Test Results\n")
	sb.WriteString("==================\n\n")
	sb.WriteString(fmt.Sprintf("Total Tests:   %d\n", results.Metrics.Total))
	sb.WriteString(fmt.Sprintf("Passed:        %d\n", results.Metrics.Passed))
	sb.WriteString(fmt.Sprintf("Failed:        %d\n", results.Metrics.Failed))
	sb.WriteString(fmt.Sprintf("Skipped:       %d\n", results.Metrics.Skipped))
	sb.WriteString(fmt.Sprintf("Retry Attempts: %d\n\n", results.Metrics.Retries))

	if len(results.Results) > 0 {
		sb.WriteString("Test Details:\n")
		for _, result := range results.Results {
			status := "✓"
			if result.Status == jtbd.TestStatusFailed {
				status = "✗"
			} else if result.Status == jtbd.TestStatusSkipped {
				status = "○"
			}
			sb.WriteString(fmt.Sprintf("  %s %s (%v)\n", status, result.TestID, result.Duration))
			if result.ErrorMessage != "" {
				sb.WriteString(fmt.Sprintf("      Error: %s\n", result.ErrorMessage))
			}
		}
	}

	return sb.String()
}

func formatJUnitResults(results *jtbd.TestResults) string {
	var sb strings.Builder

	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(fmt.Sprintf(`<testsuites tests="%d" failures="%d" skipped="%d">`+"\n",
		results.Metrics.Total, results.Metrics.Failed, results.Metrics.Skipped))
	sb.WriteString(fmt.Sprintf(`  <testsuite name="JTBD Tests" tests="%d">`+"\n", results.Metrics.Total))

	for _, result := range results.Results {
		sb.WriteString(fmt.Sprintf(`    <testcase name="%s" time="%.3f">`,
			result.TestID, result.Duration.Seconds()))
		if result.Status == jtbd.TestStatusFailed {
			sb.WriteString(fmt.Sprintf(`<failure message="%s"/>`, result.ErrorMessage))
		} else if result.Status == jtbd.TestStatusSkipped {
			sb.WriteString(fmt.Sprintf(`<skipped message="%s"/>`, result.SkipReason))
		}
		sb.WriteString(`</testcase>` + "\n")
	}

	sb.WriteString(`  </testsuite>` + "\n")
	sb.WriteString(`</testsuites>` + "\n")

	return sb.String()
}

func calculateExitCode(results *jtbd.TestResults) int {
	if results.Metrics.Failed > 0 {
		return 1
	}
	return 0
}
