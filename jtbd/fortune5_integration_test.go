// Package jtbd provides Fortune 5 JTBD integration tests.
package jtbd

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestFortune5Integration demonstrates the complete JTBD testing workflow for all Fortune 5 companies.
func TestFortune5Integration(t *testing.T) {
	t.Run("Walmart", func(t *testing.T) {
		testFortune5Company(t, "Walmart", "retail", "Stock pantry for the month")
	})

	t.Run("Amazon", func(t *testing.T) {
		testFortune5Company(t, "Amazon", "ecommerce", "Find perfect gift quickly")
	})

	t.Run("Apple", func(t *testing.T) {
		testFortune5Company(t, "Apple", "technology", "Stay connected with family")
	})

	t.Run("CVSHealth", func(t *testing.T) {
		testFortune5Company(t, "CVS Health", "healthcare", "Manage prescriptions easily")
	})

	t.Run("UnitedHealth", func(t *testing.T) {
		testFortune5Company(t, "UnitedHealth Group", "insurance", "Understand medical coverage")
	})
}

func testFortune5Company(t *testing.T, company, industry, jobName string) {
	// 1. Generate test cases using test case generator
	gen := NewTestCaseGenerator()
	testCases := gen.GenerateTestCases(industry, TestGenerationOptions{
		IncludeHappyPath:   true,
		IncludeEdgeCases:   false,
		CombinatorialLevel: 1,
	})

	if len(testCases) == 0 {
		t.Fatalf("%s: No test cases generated for %s industry", company, industry)
	}

	t.Logf("✓ %s: Generated %d test cases", company, len(testCases))

	// 2. Convert test case to Job
	testCase := testCases[0]
	job := testCase.ToJob()

	if job == nil {
		t.Fatalf("%s: Failed to convert test case to job", company)
	}

	t.Logf("✓ %s: Created JTBD job: %s", company, job.Name)

	// 3. Validate job structure
	if err := AssertJobCompleted(context.Background(), job); err != nil {
		t.Errorf("%s: Job validation failed: %v", company, err)
	}

	t.Logf("✓ %s: Job validation passed", company)

	// 4. Register job
	registry := NewJobRegistry()
	if err := registry.RegisterJob(job); err != nil {
		t.Fatalf("%s: Failed to register job: %v", company, err)
	}

	t.Logf("✓ %s: Job registered successfully", company)

	// 5. Verify job retrieval
	retrievedJob, err := registry.GetJob(job.ID)
	if err != nil {
		t.Errorf("%s: Failed to retrieve job: %v", company, err)
	}
	if retrievedJob.ID != job.ID {
		t.Errorf("%s: Retrieved job ID mismatch: got %s, want %s", company, retrievedJob.ID, job.ID)
	}

	t.Logf("✓ %s: Job retrieval verified", company)

	// 6. Test assertions
	testAssertions(t, company, job)

	t.Logf("✅ %s: All tests passed for %s industry", company, industry)
}

func testAssertions(t *testing.T, company string, job *Job) {
	// Test time compliance
	err := AssertTimeCompliance(100*time.Millisecond, 1*time.Second)
	if err != nil {
		t.Errorf("%s: Time compliance assertion failed: %v", company, err)
	}

	// Test constraint validation
	result := Result{
		JobID:   job.ID,
		Success: true,
		Data: map[string]interface{}{
			"duration": 0.5,
			"cost":     100.0,
		},
	}

	constraints := []AssertionConstraint{
		{Name: "duration", Type: "max", Value: 1.0},
		{Name: "cost", Type: "max", Value: 200.0},
	}

	err = AssertWithinConstraints(result, constraints)
	if err != nil {
		t.Errorf("%s: Constraint validation failed: %v", company, err)
	}

	t.Logf("✓ %s: Assertions passed", company)
}

// TestParallelTestExecution demonstrates concurrent test execution using the runner.
func TestParallelTestExecution(t *testing.T) {
	industries := []string{"retail", "ecommerce", "technology", "healthcare", "insurance"}

	var tests []*Test
	for idx, industry := range industries {
		testID := fmt.Sprintf("test-%s-%d", industry, idx)

		test := &Test{
			ID:   testID,
			Name: fmt.Sprintf("Fortune 5 %s Test", industry),
			Execute: func(ctx context.Context) error {
				// Simulate test work
				time.Sleep(10 * time.Millisecond)
				return nil
			},
			Timeout:    5 * time.Second,
			MaxRetries: 2,
		}
		tests = append(tests, test)
	}

	// Execute tests in parallel
	config := DefaultRunConfig()
	config.Mode = ExecutionModeParallel
	config.MaxWorkers = 10
	config.GlobalTimeout = 30 * time.Second

	engine, err := NewExecutionEngine(tests, config)
	if err != nil {
		t.Fatalf("Failed to create execution engine: %v", err)
	}

	results, err := engine.Run()
	if err != nil && err.Error() != "context canceled" {
		t.Errorf("Test execution failed: %v", err)
	}

	metrics := engine.GetMetrics()

	t.Logf("✅ Parallel execution completed:")
	t.Logf("   Total: %d, Passed: %d, Failed: %d, Skipped: %d",
		metrics.Total, metrics.Passed, metrics.Failed, metrics.Skipped)

	if metrics.Total != int32(len(tests)) {
		t.Errorf("Expected %d tests, got %d", len(tests), metrics.Total)
	}

	if metrics.Failed > 0 {
		for _, result := range results {
			if result.Status == TestStatusFailed {
				t.Errorf("Test %s failed: %s", result.TestID, result.ErrorMessage)
			}
		}
	}
}

// TestAssertionFramework tests the assertion framework components.
func TestAssertionFramework(t *testing.T) {
	t.Run("AssertionChain", func(t *testing.T) {
		chain := NewAssertionChain()
		chain.Add(AssertionResult{Pass: true, Message: "Test 1"})
		chain.Add(AssertionResult{Pass: true, Message: "Test 2"})
		chain.Add(AssertionResult{Pass: true, Message: "Test 3"})

		if !chain.IsValid() {
			t.Error("Chain should be valid with all passing assertions")
		}

		if len(chain.Results()) != 3 {
			t.Errorf("Expected 3 results, got %d", len(chain.Results()))
		}

		t.Logf("✓ Assertion chain: %s", chain.String())
	})

	t.Run("AssertionReport", func(t *testing.T) {
		report := NewAssertionReport()
		report.AddResult(AssertionResult{Pass: true})
		report.AddResult(AssertionResult{Pass: true})
		report.AddResult(AssertionResult{Pass: false})
		report.AddResult(AssertionResult{Pass: true})
		report.Complete()

		if report.PassRate() != 75.0 {
			t.Errorf("Expected 75%% pass rate, got %.1f%%", report.PassRate())
		}

		t.Logf("✓ Report: %s", report.Summary())
	})

	t.Run("ProgressTracking", func(t *testing.T) {
		tracker := NewProgressTracker()

		tracker.RecordProgress("step1", map[string]interface{}{
			"items": 5,
			"total": 100.0,
		})

		tracker.RecordCheckpoint("checkpoint1")

		progress, ok := tracker.GetProgress("step1")
		if !ok {
			t.Fatal("Progress not recorded")
		}

		if items, ok := progress.Values["items"].(int); !ok || items != 5 {
			t.Errorf("Expected 5 items, got %v", progress.Values["items"])
		}

		_, ok = tracker.GetCheckpoint("checkpoint1")
		if !ok {
			t.Error("Checkpoint not recorded")
		}

		t.Log("✓ Progress tracking works correctly")
	})
}

// TestDataFactory tests the data factory functionality.
func TestDataFactory(t *testing.T) {
	factory := NewDataFactory()

	if factory == nil {
		t.Fatal("Data factory should not be nil")
	}

	// Test that factory can create products
	products := factory.GetProductsByCompany(Walmart)
	if len(products) == 0 {
		t.Error("Expected products for Walmart")
	}

	t.Logf("✓ Data factory created with %d Walmart products", len(products))
}

// TestTestCaseGenerator tests the test case generator functionality.
func TestTestCaseGenerator(t *testing.T) {
	gen := NewTestCaseGenerator()

	// Test each industry
	industries := []string{"retail", "ecommerce", "technology", "healthcare", "insurance"}

	for _, industry := range industries {
		t.Run(industry, func(t *testing.T) {
			testCases := gen.GenerateTestCases(industry, TestGenerationOptions{
				IncludeHappyPath:   true,
				IncludeEdgeCases:   false,
				CombinatorialLevel: 1,
			})

			if len(testCases) == 0 {
				t.Errorf("No test cases generated for %s", industry)
			}

			// Verify each test case can be converted to a job
			for _, tc := range testCases {
				job := tc.ToJob()
				if job == nil {
					t.Errorf("Failed to convert test case to job for %s", industry)
				}
				if job.ID == "" {
					t.Error("Job ID should not be empty")
				}
			}

			t.Logf("✓ %s: Generated %d valid test cases", industry, len(testCases))
		})
	}
}
