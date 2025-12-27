// Package jtbd provides a comprehensive Jobs-to-be-Done (JTBD) testing framework
// for Fortune 5 companies. The framework enables testing product features and services
// by focusing on the "jobs" customers are trying to accomplish rather than demographics
// or product features.
//
// JTBD Theory Background:
// The Jobs-to-be-Done framework, popularized by Clayton Christensen, focuses on
// understanding the progress customers are trying to make in particular circumstances.
// A "job" is the progress that a person is trying to make in a particular circumstance.
//
// This implementation provides:
//   - Job definition across functional, emotional, and social dimensions
//   - Circumstance modeling (when, where, why the job arises)
//   - Outcome specifications and success criteria
//   - Progress measurement and tracking
//   - Industry-agnostic design suitable for Fortune 5 companies
//
// Fortune 5 Example Jobs:
//   - Walmart: "Help me stock my pantry for the month without multiple trips"
//   - Amazon: "Help me find the perfect gift quickly under time pressure"
//   - Apple: "Help me stay connected with my family seamlessly"
//   - CVS Health: "Help me manage my prescriptions without confusion"
//   - UnitedHealth: "Help me understand my medical coverage before a procedure"
//
// Design Principles:
//   - Pure Go implementation (no LLM dependencies)
//   - Concurrent-safe with proper mutex usage
//   - Testable and mockable interfaces
//   - Industry and domain agnostic
//   - Focused on customer progress, not product features
package jtbd

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ProgressIndicator defines the interface for measuring progress toward job completion.
// Implementations should provide metrics and signals that indicate whether the customer
// is making progress toward their desired outcome.
type ProgressIndicator interface {
	// Measure calculates the current progress value (typically 0.0 to 1.0)
	Measure(ctx context.Context, job *Job) (float64, error)

	// GetName returns the human-readable name of this progress indicator
	GetName() string

	// GetType returns the type of indicator (leading, lagging, or concurrent)
	GetType() IndicatorType

	// IsComplete returns true if this indicator shows the job is complete
	IsComplete(ctx context.Context, job *Job) (bool, error)
}

// JobTest defines the interface for testing whether a job is being fulfilled.
// Implementations should validate that a product/service helps customers make
// the desired progress in their specific circumstances.
type JobTest interface {
	// Execute runs the test and returns whether the job is being fulfilled
	Execute(ctx context.Context, job *Job) (*TestResult, error)

	// GetTestName returns the name of this test
	GetTestName() string

	// GetDescription returns a description of what this test validates
	GetDescription() string

	// Validate checks if the test configuration is valid
	Validate() error
}

// JobDimension represents the three dimensions of a job according to JTBD theory
type JobDimension string

const (
	// DimensionFunctional represents the practical, task-oriented aspect of the job
	// Example: "Get groceries into my house"
	DimensionFunctional JobDimension = "functional"

	// DimensionEmotional represents how the customer wants to feel
	// Example: "Feel confident I have enough food for my family"
	DimensionEmotional JobDimension = "emotional"

	// DimensionSocial represents how the customer wants to be perceived by others
	// Example: "Be seen as a responsible provider by my family"
	DimensionSocial JobDimension = "social"
)

// IndicatorType classifies when the indicator provides information
type IndicatorType string

const (
	// IndicatorTypeLeading measures activities that predict future success
	IndicatorTypeLeading IndicatorType = "leading"

	// IndicatorTypeLagging measures final outcomes after job completion
	IndicatorTypeLagging IndicatorType = "lagging"

	// IndicatorTypeConcurrent measures real-time progress during job execution
	IndicatorTypeConcurrent IndicatorType = "concurrent"
)

// OutcomeType categorizes the type of outcome being measured
type OutcomeType string

const (
	// OutcomeTypeSpeed focuses on how quickly the job can be completed
	OutcomeTypeSpeed OutcomeType = "speed"

	// OutcomeTypeQuality focuses on how well the job is completed
	OutcomeTypeQuality OutcomeType = "quality"

	// OutcomeTypeCost focuses on the resources required to complete the job
	OutcomeTypeCost OutcomeType = "cost"

	// OutcomeTypeExperience focuses on the customer's experience during the job
	OutcomeTypeExperience OutcomeType = "experience"
)

// CircumstanceType categorizes the context in which a job arises
type CircumstanceType string

const (
	// CircumstanceTypeTemporal relates to timing (when the job needs to be done)
	CircumstanceTypeTemporal CircumstanceType = "temporal"

	// CircumstanceTypeSpatial relates to location (where the job is being done)
	CircumstanceTypeSpatial CircumstanceType = "spatial"

	// CircumstanceTypeSituational relates to the broader situation or trigger
	CircumstanceTypeSituational CircumstanceType = "situational"

	// CircumstanceTypeSocial relates to social context (who is involved/watching)
	CircumstanceTypeSocial CircumstanceType = "social"
)

// Job represents a complete Jobs-to-be-Done definition including all three dimensions
// (functional, emotional, social), the circumstances in which it arises, and the
// desired outcomes.
//
// Example - Walmart grocery shopping:
//
//	job := &Job{
//	    ID:          "walmart-monthly-pantry-stock",
//	    Name:        "Stock pantry for the month",
//	    Description: "Get all necessary groceries to feed family for a month",
//	    Functional:  "Purchase and transport one month of groceries",
//	    Emotional:   "Feel confident my family won't run out of food",
//	    Social:      "Be seen as a reliable provider by my family",
//	    Industry:    "retail",
//	    Company:     "walmart",
//	}
type Job struct {
	// ID is a unique identifier for this job definition
	ID string

	// Name is a concise, customer-centric job statement
	Name string

	// Description provides additional context about the job
	Description string

	// Functional describes the practical task the customer wants to accomplish
	// Format: verb + object + clarifier (e.g., "Get groceries for a month")
	Functional string

	// Emotional describes how the customer wants to feel when the job is done
	// Format: "Feel [emotion] about [aspect]" (e.g., "Feel confident about food availability")
	Emotional string

	// Social describes how the customer wants to be perceived by others
	// Format: "Be seen as [identity] by [audience]" (e.g., "Be seen as organized by colleagues")
	Social string

	// Circumstances contains all the contextual factors that trigger or shape this job
	Circumstances []*Circumstance

	// Outcomes contains all the desired outcomes the customer wants to achieve
	Outcomes []*Outcome

	// Indicators contains all progress indicators for measuring job completion
	Indicators []ProgressIndicator

	// Industry is the industry sector (retail, healthcare, technology, etc.)
	Industry string

	// Company is the specific company context (walmart, amazon, apple, etc.)
	Company string

	// Metadata contains additional custom properties
	Metadata map[string]interface{}

	// CreatedAt is when this job definition was created
	CreatedAt time.Time

	// UpdatedAt is when this job definition was last modified
	UpdatedAt time.Time

	// mu protects concurrent access to job fields
	mu sync.RWMutex
}

// Circumstance represents the context in which a job arises. According to JTBD theory,
// circumstances are more important than customer demographics in predicting behavior.
//
// Example - Time pressure for gift shopping:
//
//	circ := &Circumstance{
//	    Type:        CircumstanceTypeTemporal,
//	    Description: "Only 2 days before birthday",
//	    Constraints: map[string]interface{}{
//	        "time_available": "2 days",
//	        "urgency_level": "high",
//	    },
//	}
type Circumstance struct {
	// Type categorizes this circumstance (temporal, spatial, situational, social)
	Type CircumstanceType

	// Description is a human-readable description of this circumstance
	Description string

	// Constraints contains measurable aspects of this circumstance
	// Example: {"time_available": "30 minutes", "budget_limit": 100.00}
	Constraints map[string]interface{}

	// Triggers describes what causes this circumstance to arise
	// Example: "Birthday notification received 2 days in advance"
	Triggers []string

	// Intensity represents how strongly this circumstance affects the job (0.0 to 1.0)
	// Higher values indicate more pressing or constraining circumstances
	Intensity float64

	// Metadata contains additional custom properties
	Metadata map[string]interface{}
}

// Outcome represents a desired result that indicates job completion. JTBD theory
// emphasizes that customers "hire" products to achieve specific outcomes, not features.
//
// Example - Fast prescription refill:
//
//	outcome := &Outcome{
//	    Type:        OutcomeTypeSpeed,
//	    Description: "Get prescription refilled in under 5 minutes",
//	    Metric:      "refill_time_seconds",
//	    Target:      300.0, // 5 minutes
//	    Unit:        "seconds",
//	    Priority:    1,
//	}
type Outcome struct {
	// Type categorizes this outcome (speed, quality, cost, experience)
	Type OutcomeType

	// Description is a human-readable description of the desired outcome
	Description string

	// Metric is the measurable variable that tracks this outcome
	// Example: "time_to_checkout", "error_rate", "satisfaction_score"
	Metric string

	// Target is the desired value for the metric
	Target float64

	// Unit is the unit of measurement for the metric
	// Example: "seconds", "dollars", "errors_per_100_attempts"
	Unit string

	// Priority indicates the relative importance of this outcome (1 = highest)
	Priority int

	// Direction specifies whether higher or lower values are better
	// "minimize" for metrics like time/cost, "maximize" for metrics like satisfaction
	Direction string

	// Threshold is the minimum acceptable value for this outcome
	// If actual value doesn't meet threshold, the job is not considered complete
	Threshold float64

	// Metadata contains additional custom properties
	Metadata map[string]interface{}
}

// JobRegistry manages a collection of job definitions and provides concurrent-safe
// access to job data. This is the central repository for all JTBD definitions
// in the testing framework.
type JobRegistry struct {
	mu          sync.RWMutex
	jobs        map[string]*Job
	jobsByIndustry map[string][]*Job
	jobsByCompany  map[string][]*Job
}

// NewJobRegistry creates a new JobRegistry instance
func NewJobRegistry() *JobRegistry {
	return &JobRegistry{
		jobs:           make(map[string]*Job),
		jobsByIndustry: make(map[string][]*Job),
		jobsByCompany:  make(map[string][]*Job),
	}
}

// RegisterJob adds a job to the registry
func (jr *JobRegistry) RegisterJob(job *Job) error {
	if job == nil {
		return NewJTBDError(ErrCodeInvalidJob, "job cannot be nil", nil)
	}
	if job.ID == "" {
		return NewJTBDError(ErrCodeInvalidJob, "job ID cannot be empty", nil)
	}
	if job.Name == "" {
		return NewJTBDError(ErrCodeInvalidJob, "job name cannot be empty", nil)
	}

	jr.mu.Lock()
	defer jr.mu.Unlock()

	// Set timestamps
	now := time.Now()
	if job.CreatedAt.IsZero() {
		job.CreatedAt = now
	}
	job.UpdatedAt = now

	// Initialize metadata if nil
	if job.Metadata == nil {
		job.Metadata = make(map[string]interface{})
	}

	// Store in main registry
	jr.jobs[job.ID] = job

	// Index by industry
	if job.Industry != "" {
		jr.jobsByIndustry[job.Industry] = append(jr.jobsByIndustry[job.Industry], job)
	}

	// Index by company
	if job.Company != "" {
		jr.jobsByCompany[job.Company] = append(jr.jobsByCompany[job.Company], job)
	}

	return nil
}

// GetJob retrieves a job by ID
func (jr *JobRegistry) GetJob(id string) (*Job, error) {
	jr.mu.RLock()
	defer jr.mu.RUnlock()

	job, exists := jr.jobs[id]
	if !exists {
		return nil, NewJTBDError(ErrCodeJobNotFound, fmt.Sprintf("job %q not found", id), nil)
	}
	return job, nil
}

// ListJobs returns all registered jobs
func (jr *JobRegistry) ListJobs() []*Job {
	jr.mu.RLock()
	defer jr.mu.RUnlock()

	jobs := make([]*Job, 0, len(jr.jobs))
	for _, job := range jr.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// ListJobsByIndustry returns all jobs for a specific industry
func (jr *JobRegistry) ListJobsByIndustry(industry string) []*Job {
	jr.mu.RLock()
	defer jr.mu.RUnlock()

	jobs, exists := jr.jobsByIndustry[industry]
	if !exists {
		return []*Job{}
	}
	return jobs
}

// ListJobsByCompany returns all jobs for a specific company
func (jr *JobRegistry) ListJobsByCompany(company string) []*Job {
	jr.mu.RLock()
	defer jr.mu.RUnlock()

	jobs, exists := jr.jobsByCompany[company]
	if !exists {
		return []*Job{}
	}
	return jobs
}

// RemoveJob removes a job from the registry
func (jr *JobRegistry) RemoveJob(id string) error {
	jr.mu.Lock()
	defer jr.mu.Unlock()

	job, exists := jr.jobs[id]
	if !exists {
		return NewJTBDError(ErrCodeJobNotFound, fmt.Sprintf("job %q not found", id), nil)
	}

	delete(jr.jobs, id)

	// Remove from industry index
	if job.Industry != "" {
		jr.removeFromSlice(jr.jobsByIndustry[job.Industry], id)
	}

	// Remove from company index
	if job.Company != "" {
		jr.removeFromSlice(jr.jobsByCompany[job.Company], id)
	}

	return nil
}

// removeFromSlice is a helper to remove a job from an index slice
func (jr *JobRegistry) removeFromSlice(jobs []*Job, id string) []*Job {
	for i, job := range jobs {
		if job.ID == id {
			return append(jobs[:i], jobs[i+1:]...)
		}
	}
	return jobs
}

// TestResult represents the outcome of a job test execution
type TestResult struct {
	// TestName is the name of the test that was executed
	TestName string

	// JobID is the ID of the job being tested
	JobID string

	// Success indicates whether the test passed
	Success bool

	// Score is a numeric score (typically 0.0 to 1.0) indicating how well the job was fulfilled
	Score float64

	// Message provides human-readable details about the test result
	Message string

	// ProgressMeasurements contains individual progress indicator measurements
	ProgressMeasurements map[string]float64

	// OutcomeResults contains results for each outcome
	OutcomeResults map[string]*OutcomeResult

	// ExecutionTime is how long the test took to execute
	ExecutionTime time.Duration

	// Timestamp is when the test was executed
	Timestamp time.Time

	// Metadata contains additional custom properties
	Metadata map[string]interface{}
}

// OutcomeResult represents the result of measuring a specific outcome
type OutcomeResult struct {
	// OutcomeDescription is the description of the outcome being measured
	OutcomeDescription string

	// MetricName is the name of the metric being measured
	MetricName string

	// ActualValue is the measured value
	ActualValue float64

	// TargetValue is the desired value
	TargetValue float64

	// ThresholdValue is the minimum acceptable value
	ThresholdValue float64

	// Unit is the unit of measurement
	Unit string

	// MetThreshold indicates if the threshold was met
	MetThreshold bool

	// MetTarget indicates if the target was met
	MetTarget bool

	// PerformanceRatio is ActualValue / TargetValue (adjusted for direction)
	PerformanceRatio float64
}

// TestExecutor manages the execution of job tests
type TestExecutor struct {
	mu       sync.RWMutex
	registry *JobRegistry
	tests    map[string]JobTest
	results  []*TestResult
}

// NewTestExecutor creates a new TestExecutor instance
func NewTestExecutor(registry *JobRegistry) *TestExecutor {
	return &TestExecutor{
		registry: registry,
		tests:    make(map[string]JobTest),
		results:  make([]*TestResult, 0),
	}
}

// RegisterTest adds a test to the executor
func (te *TestExecutor) RegisterTest(test JobTest) error {
	if test == nil {
		return NewJTBDError(ErrCodeInvalidTest, "test cannot be nil", nil)
	}

	if err := test.Validate(); err != nil {
		return NewJTBDError(ErrCodeInvalidTest, "test validation failed", err)
	}

	te.mu.Lock()
	defer te.mu.Unlock()

	te.tests[test.GetTestName()] = test
	return nil
}

// ExecuteTest runs a specific test against a job
func (te *TestExecutor) ExecuteTest(ctx context.Context, testName string, jobID string) (*TestResult, error) {
	te.mu.RLock()
	test, exists := te.tests[testName]
	te.mu.RUnlock()

	if !exists {
		return nil, NewJTBDError(ErrCodeTestNotFound, fmt.Sprintf("test %q not found", testName), nil)
	}

	job, err := te.registry.GetJob(jobID)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	result, err := test.Execute(ctx, job)
	if err != nil {
		return nil, err
	}

	result.ExecutionTime = time.Since(startTime)
	result.Timestamp = time.Now()

	// Store result
	te.mu.Lock()
	te.results = append(te.results, result)
	te.mu.Unlock()

	return result, nil
}

// ExecuteAllTests runs all registered tests against a job
func (te *TestExecutor) ExecuteAllTests(ctx context.Context, jobID string) ([]*TestResult, error) {
	te.mu.RLock()
	testNames := make([]string, 0, len(te.tests))
	for name := range te.tests {
		testNames = append(testNames, name)
	}
	te.mu.RUnlock()

	results := make([]*TestResult, 0, len(testNames))
	for _, testName := range testNames {
		result, err := te.ExecuteTest(ctx, testName, jobID)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// GetResults returns all test results
func (te *TestExecutor) GetResults() []*TestResult {
	te.mu.RLock()
	defer te.mu.RUnlock()

	// Return a copy to prevent external modification
	results := make([]*TestResult, len(te.results))
	copy(results, te.results)
	return results
}

// ClearResults removes all stored test results
func (te *TestExecutor) ClearResults() {
	te.mu.Lock()
	defer te.mu.Unlock()

	te.results = make([]*TestResult, 0)
}

// JobBuilder provides a fluent API for constructing Job definitions
type JobBuilder struct {
	job *Job
	err error
}

// NewJobBuilder creates a new JobBuilder instance
func NewJobBuilder(id, name string) *JobBuilder {
	return &JobBuilder{
		job: &Job{
			ID:            id,
			Name:          name,
			Circumstances: make([]*Circumstance, 0),
			Outcomes:      make([]*Outcome, 0),
			Indicators:    make([]ProgressIndicator, 0),
			Metadata:      make(map[string]interface{}),
		},
	}
}

// WithDescription sets the job description
func (jb *JobBuilder) WithDescription(description string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Description = description
	return jb
}

// WithFunctional sets the functional dimension
func (jb *JobBuilder) WithFunctional(functional string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Functional = functional
	return jb
}

// WithEmotional sets the emotional dimension
func (jb *JobBuilder) WithEmotional(emotional string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Emotional = emotional
	return jb
}

// WithSocial sets the social dimension
func (jb *JobBuilder) WithSocial(social string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Social = social
	return jb
}

// WithIndustry sets the industry
func (jb *JobBuilder) WithIndustry(industry string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Industry = industry
	return jb
}

// WithCompany sets the company
func (jb *JobBuilder) WithCompany(company string) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Company = company
	return jb
}

// AddCircumstance adds a circumstance to the job
func (jb *JobBuilder) AddCircumstance(circumstance *Circumstance) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Circumstances = append(jb.job.Circumstances, circumstance)
	return jb
}

// AddOutcome adds an outcome to the job
func (jb *JobBuilder) AddOutcome(outcome *Outcome) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Outcomes = append(jb.job.Outcomes, outcome)
	return jb
}

// AddIndicator adds a progress indicator to the job
func (jb *JobBuilder) AddIndicator(indicator ProgressIndicator) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Indicators = append(jb.job.Indicators, indicator)
	return jb
}

// WithMetadata sets metadata fields
func (jb *JobBuilder) WithMetadata(key string, value interface{}) *JobBuilder {
	if jb.err != nil {
		return jb
	}
	jb.job.Metadata[key] = value
	return jb
}

// Build finalizes the job construction and returns the job
func (jb *JobBuilder) Build() (*Job, error) {
	if jb.err != nil {
		return nil, jb.err
	}

	// Validate required fields
	if jb.job.ID == "" {
		return nil, NewJTBDError(ErrCodeInvalidJob, "job ID is required", nil)
	}
	if jb.job.Name == "" {
		return nil, NewJTBDError(ErrCodeInvalidJob, "job name is required", nil)
	}

	return jb.job, nil
}

// JTBDError represents errors specific to the JTBD framework
type JTBDError struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface
func (je *JTBDError) Error() string {
	if je.Cause != nil {
		return fmt.Sprintf("jtbd: %s (%s): %v", je.Code, je.Message, je.Cause)
	}
	return fmt.Sprintf("jtbd: %s (%s)", je.Code, je.Message)
}

// Unwrap returns the underlying cause error
func (je *JTBDError) Unwrap() error {
	return je.Cause
}

// NewJTBDError creates a new JTBDError
func NewJTBDError(code, message string, cause error) *JTBDError {
	return &JTBDError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Common error codes
const (
	ErrCodeInvalidJob    = "invalid_job"
	ErrCodeInvalidTest   = "invalid_test"
	ErrCodeJobNotFound   = "job_not_found"
	ErrCodeTestNotFound  = "test_not_found"
	ErrCodeTestFailed    = "test_failed"
	ErrCodeInvalidInput  = "invalid_input"
	ErrCodeInternalError = "internal_error"
)
