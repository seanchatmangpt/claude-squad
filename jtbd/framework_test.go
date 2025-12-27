package jtbd

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestJobRegistry_RegisterJob(t *testing.T) {
	registry := NewJobRegistry()

	job := &Job{
		ID:          "test-job-1",
		Name:        "Test Job",
		Description: "A test job for unit testing",
		Functional:  "Complete a test",
		Emotional:   "Feel confident about testing",
		Social:      "Be seen as thorough by the team",
		Industry:    "technology",
		Company:     "test-company",
	}

	err := registry.RegisterJob(job)
	if err != nil {
		t.Fatalf("Failed to register job: %v", err)
	}

	// Verify job was registered
	retrieved, err := registry.GetJob("test-job-1")
	if err != nil {
		t.Fatalf("Failed to retrieve job: %v", err)
	}

	if retrieved.ID != job.ID {
		t.Errorf("Expected job ID %s, got %s", job.ID, retrieved.ID)
	}

	if retrieved.Name != job.Name {
		t.Errorf("Expected job name %s, got %s", job.Name, retrieved.Name)
	}
}

func TestJobRegistry_RegisterJob_NilJob(t *testing.T) {
	registry := NewJobRegistry()

	err := registry.RegisterJob(nil)
	if err == nil {
		t.Fatal("Expected error when registering nil job, got nil")
	}

	jtbdErr, ok := err.(*JTBDError)
	if !ok {
		t.Fatalf("Expected JTBDError, got %T", err)
	}

	if jtbdErr.Code != ErrCodeInvalidJob {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidJob, jtbdErr.Code)
	}
}

func TestJobRegistry_RegisterJob_EmptyID(t *testing.T) {
	registry := NewJobRegistry()

	job := &Job{
		ID:   "", // Empty ID
		Name: "Test Job",
	}

	err := registry.RegisterJob(job)
	if err == nil {
		t.Fatal("Expected error when registering job with empty ID, got nil")
	}
}

func TestJobRegistry_GetJob_NotFound(t *testing.T) {
	registry := NewJobRegistry()

	_, err := registry.GetJob("non-existent-job")
	if err == nil {
		t.Fatal("Expected error when getting non-existent job, got nil")
	}

	jtbdErr, ok := err.(*JTBDError)
	if !ok {
		t.Fatalf("Expected JTBDError, got %T", err)
	}

	if jtbdErr.Code != ErrCodeJobNotFound {
		t.Errorf("Expected error code %s, got %s", ErrCodeJobNotFound, jtbdErr.Code)
	}
}

func TestJobRegistry_ListJobs(t *testing.T) {
	registry := NewJobRegistry()

	// Register multiple jobs
	jobs := []*Job{
		{ID: "job-1", Name: "Job 1", Industry: "retail"},
		{ID: "job-2", Name: "Job 2", Industry: "healthcare"},
		{ID: "job-3", Name: "Job 3", Industry: "retail"},
	}

	for _, job := range jobs {
		if err := registry.RegisterJob(job); err != nil {
			t.Fatalf("Failed to register job: %v", err)
		}
	}

	// List all jobs
	allJobs := registry.ListJobs()
	if len(allJobs) != 3 {
		t.Errorf("Expected 3 jobs, got %d", len(allJobs))
	}
}

func TestJobRegistry_ListJobsByIndustry(t *testing.T) {
	registry := NewJobRegistry()

	// Register jobs in different industries
	jobs := []*Job{
		{ID: "job-1", Name: "Job 1", Industry: "retail"},
		{ID: "job-2", Name: "Job 2", Industry: "healthcare"},
		{ID: "job-3", Name: "Job 3", Industry: "retail"},
	}

	for _, job := range jobs {
		if err := registry.RegisterJob(job); err != nil {
			t.Fatalf("Failed to register job: %v", err)
		}
	}

	// List jobs by industry
	retailJobs := registry.ListJobsByIndustry("retail")
	if len(retailJobs) != 2 {
		t.Errorf("Expected 2 retail jobs, got %d", len(retailJobs))
	}

	healthcareJobs := registry.ListJobsByIndustry("healthcare")
	if len(healthcareJobs) != 1 {
		t.Errorf("Expected 1 healthcare job, got %d", len(healthcareJobs))
	}

	// Non-existent industry
	techJobs := registry.ListJobsByIndustry("technology")
	if len(techJobs) != 0 {
		t.Errorf("Expected 0 technology jobs, got %d", len(techJobs))
	}
}

func TestJobRegistry_RemoveJob(t *testing.T) {
	registry := NewJobRegistry()

	job := &Job{
		ID:       "job-to-remove",
		Name:     "Job to Remove",
		Industry: "retail",
	}

	if err := registry.RegisterJob(job); err != nil {
		t.Fatalf("Failed to register job: %v", err)
	}

	// Remove the job
	if err := registry.RemoveJob("job-to-remove"); err != nil {
		t.Fatalf("Failed to remove job: %v", err)
	}

	// Verify job was removed
	_, err := registry.GetJob("job-to-remove")
	if err == nil {
		t.Fatal("Expected error when getting removed job, got nil")
	}
}

func TestJobBuilder(t *testing.T) {
	job, err := NewJobBuilder("test-job", "Test Job").
		WithDescription("A test job built with builder pattern").
		WithFunctional("Complete the task").
		WithEmotional("Feel accomplished").
		WithSocial("Be recognized by peers").
		WithIndustry("technology").
		WithCompany("test-company").
		WithMetadata("key1", "value1").
		AddCircumstance(&Circumstance{
			Type:        CircumstanceTypeTemporal,
			Description: "Under time pressure",
			Intensity:   0.8,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeSpeed,
			Description: "Complete quickly",
			Metric:      "completion_time",
			Target:      60.0,
			Unit:        "seconds",
			Priority:    1,
		}).
		Build()

	if err != nil {
		t.Fatalf("Failed to build job: %v", err)
	}

	if job.ID != "test-job" {
		t.Errorf("Expected ID 'test-job', got %s", job.ID)
	}

	if job.Name != "Test Job" {
		t.Errorf("Expected name 'Test Job', got %s", job.Name)
	}

	if len(job.Circumstances) != 1 {
		t.Errorf("Expected 1 circumstance, got %d", len(job.Circumstances))
	}

	if len(job.Outcomes) != 1 {
		t.Errorf("Expected 1 outcome, got %d", len(job.Outcomes))
	}

	if job.Metadata["key1"] != "value1" {
		t.Errorf("Expected metadata key1='value1', got %v", job.Metadata["key1"])
	}
}

func TestJobBuilder_MissingRequiredFields(t *testing.T) {
	// Test missing ID
	_, err := NewJobBuilder("", "Test Job").Build()
	if err == nil {
		t.Fatal("Expected error when building job with empty ID, got nil")
	}

	// Test missing name
	_, err = NewJobBuilder("test-id", "").Build()
	if err == nil {
		t.Fatal("Expected error when building job with empty name, got nil")
	}
}

func TestTestExecutor_RegisterTest(t *testing.T) {
	registry := NewJobRegistry()
	executor := NewTestExecutor(registry)

	test := NewSimpleJobTest(
		"test-1",
		"A simple test",
		func(ctx context.Context, job *Job) (*TestResult, error) {
			return &TestResult{
				TestName: "test-1",
				JobID:    job.ID,
				Success:  true,
				Score:    1.0,
			}, nil
		},
	)

	err := executor.RegisterTest(test)
	if err != nil {
		t.Fatalf("Failed to register test: %v", err)
	}
}

func TestTestExecutor_RegisterTest_NilTest(t *testing.T) {
	registry := NewJobRegistry()
	executor := NewTestExecutor(registry)

	err := executor.RegisterTest(nil)
	if err == nil {
		t.Fatal("Expected error when registering nil test, got nil")
	}
}

func TestTestExecutor_ExecuteTest(t *testing.T) {
	registry := NewJobRegistry()
	executor := NewTestExecutor(registry)

	// Register a job
	job := &Job{
		ID:   "test-job",
		Name: "Test Job",
	}
	if err := registry.RegisterJob(job); err != nil {
		t.Fatalf("Failed to register job: %v", err)
	}

	// Register a test
	test := NewSimpleJobTest(
		"test-1",
		"A simple test",
		func(ctx context.Context, j *Job) (*TestResult, error) {
			return &TestResult{
				TestName: "test-1",
				JobID:    j.ID,
				Success:  true,
				Score:    0.95,
				Message:  "Test passed successfully",
			}, nil
		},
	)
	if err := executor.RegisterTest(test); err != nil {
		t.Fatalf("Failed to register test: %v", err)
	}

	// Execute the test
	ctx := context.Background()
	result, err := executor.ExecuteTest(ctx, "test-1", "test-job")
	if err != nil {
		t.Fatalf("Failed to execute test: %v", err)
	}

	if !result.Success {
		t.Error("Expected test to succeed")
	}

	if result.Score != 0.95 {
		t.Errorf("Expected score 0.95, got %.2f", result.Score)
	}

	if result.TestName != "test-1" {
		t.Errorf("Expected test name 'test-1', got %s", result.TestName)
	}
}

func TestTestExecutor_ExecuteTest_TestNotFound(t *testing.T) {
	registry := NewJobRegistry()
	executor := NewTestExecutor(registry)

	ctx := context.Background()
	_, err := executor.ExecuteTest(ctx, "non-existent-test", "some-job")
	if err == nil {
		t.Fatal("Expected error when executing non-existent test, got nil")
	}

	jtbdErr, ok := err.(*JTBDError)
	if !ok {
		t.Fatalf("Expected JTBDError, got %T", err)
	}

	if jtbdErr.Code != ErrCodeTestNotFound {
		t.Errorf("Expected error code %s, got %s", ErrCodeTestNotFound, jtbdErr.Code)
	}
}

func TestTestExecutor_ExecuteAllTests(t *testing.T) {
	registry := NewJobRegistry()
	executor := NewTestExecutor(registry)

	// Register a job
	job := &Job{
		ID:   "test-job",
		Name: "Test Job",
	}
	if err := registry.RegisterJob(job); err != nil {
		t.Fatalf("Failed to register job: %v", err)
	}

	// Register multiple tests
	for i := 1; i <= 3; i++ {
		testName := fmt.Sprintf("test-%d", i)
		test := NewSimpleJobTest(
			testName,
			fmt.Sprintf("Test number %d", i),
			func(ctx context.Context, j *Job) (*TestResult, error) {
				return &TestResult{
					TestName: testName,
					JobID:    j.ID,
					Success:  true,
					Score:    1.0,
				}, nil
			},
		)
		if err := executor.RegisterTest(test); err != nil {
			t.Fatalf("Failed to register test %s: %v", testName, err)
		}
	}

	// Execute all tests
	ctx := context.Background()
	results, err := executor.ExecuteAllTests(ctx, "test-job")
	if err != nil {
		t.Fatalf("Failed to execute all tests: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	for _, result := range results {
		if !result.Success {
			t.Errorf("Test %s failed unexpectedly", result.TestName)
		}
	}
}

func TestSimpleProgressIndicator(t *testing.T) {
	indicator := NewSimpleProgressIndicator(
		"test-indicator",
		IndicatorTypeLeading,
		func(ctx context.Context, job *Job) (float64, error) {
			return 0.75, nil
		},
	)

	if indicator.GetName() != "test-indicator" {
		t.Errorf("Expected name 'test-indicator', got %s", indicator.GetName())
	}

	if indicator.GetType() != IndicatorTypeLeading {
		t.Errorf("Expected type %s, got %s", IndicatorTypeLeading, indicator.GetType())
	}

	ctx := context.Background()
	job := &Job{ID: "test-job", Name: "Test Job"}

	progress, err := indicator.Measure(ctx, job)
	if err != nil {
		t.Fatalf("Failed to measure progress: %v", err)
	}

	if progress != 0.75 {
		t.Errorf("Expected progress 0.75, got %.2f", progress)
	}

	complete, err := indicator.IsComplete(ctx, job)
	if err != nil {
		t.Fatalf("Failed to check completion: %v", err)
	}

	if complete {
		t.Error("Expected indicator to not be complete at 0.75 progress")
	}
}

func TestSimpleProgressIndicator_IsComplete(t *testing.T) {
	indicator := NewSimpleProgressIndicator(
		"complete-indicator",
		IndicatorTypeLagging,
		func(ctx context.Context, job *Job) (float64, error) {
			return 1.0, nil
		},
	)

	ctx := context.Background()
	job := &Job{ID: "test-job", Name: "Test Job"}

	complete, err := indicator.IsComplete(ctx, job)
	if err != nil {
		t.Fatalf("Failed to check completion: %v", err)
	}

	if !complete {
		t.Error("Expected indicator to be complete at 1.0 progress")
	}
}

func TestJTBDError(t *testing.T) {
	// Test error without cause
	err1 := NewJTBDError(ErrCodeInvalidJob, "test error", nil)
	expected1 := "jtbd: invalid_job (test error)"
	if err1.Error() != expected1 {
		t.Errorf("Expected error message %q, got %q", expected1, err1.Error())
	}

	// Test error with cause
	cause := fmt.Errorf("underlying error")
	err2 := NewJTBDError(ErrCodeInvalidJob, "test error", cause)
	expected2 := "jtbd: invalid_job (test error): underlying error"
	if err2.Error() != expected2 {
		t.Errorf("Expected error message %q, got %q", expected2, err2.Error())
	}

	// Test Unwrap
	if err2.Unwrap() != cause {
		t.Error("Unwrap did not return the correct cause")
	}
}

func TestConcurrentAccess(t *testing.T) {
	registry := NewJobRegistry()
	done := make(chan bool)
	numGoroutines := 10

	// Test concurrent registration
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			job := &Job{
				ID:   fmt.Sprintf("job-%d", id),
				Name: fmt.Sprintf("Job %d", id),
			}
			if err := registry.RegisterJob(job); err != nil {
				t.Errorf("Failed to register job in goroutine: %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all jobs were registered
	jobs := registry.ListJobs()
	if len(jobs) != numGoroutines {
		t.Errorf("Expected %d jobs, got %d", numGoroutines, len(jobs))
	}
}

func TestTestExecutor_ClearResults(t *testing.T) {
	registry := NewJobRegistry()
	executor := NewTestExecutor(registry)

	// Register a job
	job := &Job{
		ID:   "test-job",
		Name: "Test Job",
	}
	if err := registry.RegisterJob(job); err != nil {
		t.Fatalf("Failed to register job: %v", err)
	}

	// Register and execute a test
	test := NewSimpleJobTest(
		"test-1",
		"A simple test",
		func(ctx context.Context, j *Job) (*TestResult, error) {
			return &TestResult{
				TestName: "test-1",
				JobID:    j.ID,
				Success:  true,
			}, nil
		},
	)
	if err := executor.RegisterTest(test); err != nil {
		t.Fatalf("Failed to register test: %v", err)
	}

	ctx := context.Background()
	if _, err := executor.ExecuteTest(ctx, "test-1", "test-job"); err != nil {
		t.Fatalf("Failed to execute test: %v", err)
	}

	// Verify results exist
	results := executor.GetResults()
	if len(results) == 0 {
		t.Fatal("Expected at least one result before clearing")
	}

	// Clear results
	executor.ClearResults()

	// Verify results are cleared
	results = executor.GetResults()
	if len(results) != 0 {
		t.Errorf("Expected 0 results after clearing, got %d", len(results))
	}
}

// Benchmark tests
func BenchmarkJobRegistry_RegisterJob(b *testing.B) {
	registry := NewJobRegistry()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		job := &Job{
			ID:   fmt.Sprintf("job-%d", i),
			Name: fmt.Sprintf("Job %d", i),
		}
		_ = registry.RegisterJob(job)
	}
}

func BenchmarkJobRegistry_GetJob(b *testing.B) {
	registry := NewJobRegistry()

	// Pre-populate registry
	for i := 0; i < 1000; i++ {
		job := &Job{
			ID:   fmt.Sprintf("job-%d", i),
			Name: fmt.Sprintf("Job %d", i),
		}
		_ = registry.RegisterJob(job)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = registry.GetJob(fmt.Sprintf("job-%d", i%1000))
	}
}

func BenchmarkTestExecutor_ExecuteTest(b *testing.B) {
	registry := NewJobRegistry()
	executor := NewTestExecutor(registry)

	job := &Job{
		ID:   "benchmark-job",
		Name: "Benchmark Job",
	}
	_ = registry.RegisterJob(job)

	test := NewSimpleJobTest(
		"benchmark-test",
		"A benchmark test",
		func(ctx context.Context, j *Job) (*TestResult, error) {
			return &TestResult{
				TestName:  "benchmark-test",
				JobID:     j.ID,
				Success:   true,
				Score:     1.0,
				Timestamp: time.Now(),
			}, nil
		},
	)
	_ = executor.RegisterTest(test)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = executor.ExecuteTest(ctx, "benchmark-test", "benchmark-job")
	}
}
