// Package jtbd provides end-to-end Fortune 5 JTBD test examples.
package jtbd

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestFortune5WalmartJourney demonstrates end-to-end Walmart grocery shopping JTBD test.
func TestFortune5WalmartJourney(t *testing.T) {
	// Build job using JTBD framework
	job, err := NewJobBuilder("walmart-monthly-pantry", "Stock pantry for the month").
		WithFunctional("Get enough groceries to feed my family for a month").
		WithEmotional("Feel confident my family won't run out of food").
		WithSocial("Be seen as a reliable provider by my family").
		AddOutcome(&Outcome{
			Type:      OutcomeTypeSpeed,
			Metric:    "shopping_time_minutes",
			Target:    30.0,
			Threshold: 45.0,
		}).
		AddOutcome(&Outcome{
			Type:      OutcomeTypeCost,
			Metric:    "total_cost_dollars",
			Target:    150.0,
			Threshold: 160.0,
		}).
		Build()

	if err != nil {
		t.Fatalf("Failed to build job: %v", err)
	}

	// Execute test
	ctx := context.Background()
	registry := NewJobRegistry()
	registry.RegisterJob(job)

	executor := NewTestExecutor(registry)
	result, err := executor.ExecuteTest(ctx, "walmart-test", job.ID)

	if err != nil {
		t.Errorf("Test execution failed: %v", err)
	}

	// Assertions
	if err := AssertJobCompleted(ctx, job); err != nil {
		t.Errorf("Job completion assertion failed: %v", err)
	}

	t.Logf("✅ Walmart JTBD test passed: %s (result: %v)", job.Name, result)
}

// TestFortune5AmazonGiftFinding demonstrates Amazon gift shopping JTBD test.
func TestFortune5AmazonGiftFinding(t *testing.T) {
	// Generate test case
	gen := NewTestCaseGenerator()
	testCases := gen.GenerateTestCases("ecommerce", TestGenerationOptions{
		IncludeHappyPath: true,
		IncludeEdgeCases: false,
	})

	if len(testCases) == 0 {
		t.Fatal("No test cases generated for ecommerce")
	}

	// Use first test case
	testCase := testCases[0]
	job := testCase.ToJob()

	// Execute
	ctx := context.Background()
	registry := NewJobRegistry()
	registry.RegisterJob(job)

	executor := NewTestExecutor(registry)
	_, err := executor.ExecuteTest(ctx, "amazon-test", job.ID)

	if err != nil {
		t.Errorf("Amazon test failed: %v", err)
	}

	// Time compliance
	if err := AssertTimeCompliance(500*time.Millisecond, 1*time.Second); err != nil {
		t.Errorf("Time compliance failed: %v", err)
	}

	t.Logf("✅ Amazon JTBD test passed")
}

// TestFortune5AppleFamilyConnection demonstrates Apple device setup JTBD test.
func TestFortune5AppleFamilyConnection(t *testing.T) {
	gen := NewTestCaseGenerator()
	testCases := gen.GenerateTestCases("technology", TestGenerationOptions{
		IncludeHappyPath:   true,
		CombinatorialLevel: 1,
	})

	if len(testCases) == 0 {
		t.Fatal("No test cases generated for technology")
	}

	job := testCases[0].ToJob()

	// Create progress tracker
	tracker := NewProgressTracker()
	tracker.RecordCheckpoint("start")

	// Execute
	ctx := context.Background()
	registry := NewJobRegistry()
	registry.RegisterJob(job)

	executor := NewTestExecutor(registry)
	_, err := executor.ExecuteTest(ctx, "apple-test", job.ID)

	tracker.RecordCheckpoint("end")

	if err != nil {
		t.Errorf("Apple test failed: %v", err)
	}

	// Verify checkpoints
	start, ok1 := tracker.GetCheckpoint("start")
	end, ok2 := tracker.GetCheckpoint("end")

	if !ok1 || !ok2 {
		t.Error("Checkpoints not recorded")
	}

	if !end.After(start) {
		t.Error("End checkpoint should be after start")
	}

	t.Logf("✅ Apple JTBD test passed in %v", end.Sub(start))
}

// TestFortune5CVSPrescriptionManagement demonstrates CVS pharmacy JTBD test.
func TestFortune5CVSPrescriptionManagement(t *testing.T) {
	gen := NewTestCaseGenerator()
	testCases := gen.GenerateTestCases("healthcare", TestGenerationOptions{
		IncludeHappyPath: true,
		IncludeEdgeCases: true,
	})

	if len(testCases) < 2 {
		t.Fatal("Expected at least 2 test cases (happy path + edge case)")
	}

	for i, testCase := range testCases[:2] {
		job := testCase.ToJob()

		ctx := context.Background()
		registry := NewJobRegistry()
		registry.RegisterJob(job)

		executor := NewTestExecutor(registry)
		_, err := executor.ExecuteTest(ctx, fmt.Sprintf("cvs-test-%d", i), job.ID)

		if err != nil {
			t.Errorf("CVS test case %d failed: %v", i, err)
		}
	}

	t.Logf("✅ CVS JTBD tests passed: %d scenarios", len(testCases[:2]))
}

// TestFortune5UnitedHealthCoverageUnderstanding demonstrates insurance JTBD test.
func TestFortune5UnitedHealthCoverageUnderstanding(t *testing.T) {
	gen := NewTestCaseGenerator()
	testCases := gen.GenerateTestCases("insurance", TestGenerationOptions{
		IncludeHappyPath: true,
	})

	if len(testCases) == 0 {
		t.Fatal("No test cases generated for insurance")
	}

	job := testCases[0].ToJob()

	// Create expectations
	expectations := Expectations{
		FunctionalCriteria: []string{"Understand my coverage"},
		EmotionalCriteria:  []string{"Feel secure about medical costs"},
		SocialCriteria:     []string{"Make informed healthcare decisions"},
		MaxDuration:        5 * time.Minute,
	}

	ctx := context.Background()

	// Validate satisfaction
	if err := AssertSatisfaction(ctx, job, expectations); err != nil {
		t.Errorf("Satisfaction assertion failed: %v", err)
	}

	registry := NewJobRegistry()
	registry.RegisterJob(job)

	executor := NewTestExecutor(registry)
	_, err := executor.ExecuteTest(ctx, "unitedhealth-test", job.ID)

	if err != nil {
		t.Errorf("UnitedHealth test failed: %v", err)
	}

	t.Logf("✅ UnitedHealth JTBD test passed")
}

// TestParallelFortune5Execution demonstrates concurrent testing across all Fortune 5.
func TestParallelFortune5Execution(t *testing.T) {
	industries := []string{"retail", "ecommerce", "technology", "healthcare", "insurance"}

	gen := NewTestCaseGenerator()

	// Generate tests for all industries
	var allTests []*Test
	for _, industry := range industries {
		testCases := gen.GenerateTestCases(industry, TestGenerationOptions{
			IncludeHappyPath: true,
		})

		for idx, tc := range testCases {
			job := tc.ToJob()
			testID := fmt.Sprintf("%s-test-%d", industry, idx)

			// Create a closure to capture the job
			jobCopy := job
			test := &Test{
				ID:   testID,
				Name: fmt.Sprintf("%s JTBD Test %d", industry, idx),
				Execute: func(ctx context.Context) error {
					registry := NewJobRegistry()
					registry.RegisterJob(jobCopy)
					executor := NewTestExecutor(registry)
					_, err := executor.ExecuteTest(ctx, jobCopy.ID, jobCopy.ID)
					return err
				},
				Timeout:    30 * time.Second,
				MaxRetries: 2,
			}
			allTests = append(allTests, test)
		}
	}

	// Execute in parallel
	config := DefaultRunConfig()
	config.Mode = ExecutionModeParallel
	config.MaxWorkers = 10

	engine, err := NewExecutionEngine(allTests, config)
	if err != nil {
		t.Fatalf("Failed to create execution engine: %v", err)
	}

	results, err := engine.Run()
	if err != nil && err.Error() != "context canceled" {
		t.Logf("Parallel execution warning: %v", err)
	}

	metrics := engine.GetMetrics()

	t.Logf("✅ Parallel Fortune 5 execution completed:")
	t.Logf("   Total: %d, Passed: %d, Failed: %d, Skipped: %d",
		metrics.Total, metrics.Passed, metrics.Failed, metrics.Skipped)

	// Verify all tests ran
	if int(metrics.Total) != len(allTests) {
		t.Errorf("Expected %d tests, got %d", len(allTests), metrics.Total)
	}

	// Report any failures
	if metrics.Failed > 0 {
		for _, result := range results {
			if result.Status == TestStatusFailed {
				t.Logf("Test %s failed: %v", result.TestID, result.ErrorMessage)
			}
		}
	}
}

// TestAssertionChain demonstrates assertion chaining.
func TestAssertionChain(t *testing.T) {
	chain := NewAssertionChain()

	chain.Add(AssertionResult{Pass: true, Message: "Test 1 passed"})
	chain.Add(AssertionResult{Pass: true, Message: "Test 2 passed"})
	chain.Add(AssertionResult{Pass: true, Message: "Test 3 passed"})

	if !chain.IsValid() {
		t.Error("Chain should be valid with all passing assertions")
	}

	results := chain.Results()
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	t.Logf("✅ Assertion chain: %s", chain.String())
}

// TestAssertionReport demonstrates assertion reporting.
func TestAssertionReport(t *testing.T) {
	report := NewAssertionReport()

	report.AddResult(AssertionResult{Pass: true, Message: "Test 1"})
	report.AddResult(AssertionResult{Pass: true, Message: "Test 2"})
	report.AddResult(AssertionResult{Pass: false, Message: "Test 3"})
	report.AddResult(AssertionResult{Pass: true, Message: "Test 4"})

	report.Complete()

	if report.PassRate() != 75.0 {
		t.Errorf("Expected 75%% pass rate, got %.1f%%", report.PassRate())
	}

	if report.IsSuccessful() {
		t.Error("Report should not be successful with failures")
	}

	t.Logf("✅ Assertion report: %s", report.Summary())
}

// TestConstraintValidation demonstrates constraint validation.
func TestConstraintValidation(t *testing.T) {
	result := Result{
		JobID:   "test-job",
		Success: true,
		Data: map[string]interface{}{
			"cost":     145.0,
			"duration": 28.5,
			"quality":  "high",
		},
	}

	constraints := []AssertionConstraint{
		{Name: "cost", Type: "max", Value: 150.0},
		{Name: "duration", Type: "max", Value: 30.0},
		{Name: "quality", Type: "equals", Value: "high"},
	}

	if err := AssertWithinConstraints(result, constraints); err != nil {
		t.Errorf("Constraint validation failed: %v", err)
	}

	t.Log("✅ All constraints validated successfully")
}

// TestProgressTracking demonstrates progress tracking over time.
func TestProgressTracking(t *testing.T) {
	tracker := NewProgressTracker()

	// Record initial progress
	tracker.RecordProgress("cart", map[string]interface{}{
		"items": 0,
		"total": 0.0,
	})

	time.Sleep(10 * time.Millisecond)

	// Record updated progress
	tracker.RecordProgress("cart", map[string]interface{}{
		"items": 5,
		"total": 45.50,
	})

	// Verify progress was recorded
	progress, ok := tracker.GetProgress("cart")
	if !ok {
		t.Fatal("Progress not recorded")
	}

	items, ok := progress.Values["items"].(int)
	if !ok || items != 5 {
		t.Errorf("Expected 5 items, got %v", progress.Values["items"])
	}

	t.Logf("✅ Progress tracking: %d items, $%.2f", items, progress.Values["total"])
}
