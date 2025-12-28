package jtbd

import (
	"context"
	"fmt"
	"time"
)

// Example implementations and usage patterns for the JTBD framework

// SimpleProgressIndicator is a basic implementation of ProgressIndicator for demonstration
type SimpleProgressIndicator struct {
	name         string
	indicatorType IndicatorType
	measureFunc  func(context.Context, *Job) (float64, error)
}

// NewSimpleProgressIndicator creates a new SimpleProgressIndicator
func NewSimpleProgressIndicator(name string, iType IndicatorType, measureFunc func(context.Context, *Job) (float64, error)) *SimpleProgressIndicator {
	return &SimpleProgressIndicator{
		name:         name,
		indicatorType: iType,
		measureFunc:  measureFunc,
	}
}

// Measure implements ProgressIndicator
func (spi *SimpleProgressIndicator) Measure(ctx context.Context, job *Job) (float64, error) {
	if spi.measureFunc == nil {
		return 0.0, NewJTBDError(ErrCodeInternalError, "measure function not set", nil)
	}
	return spi.measureFunc(ctx, job)
}

// GetName implements ProgressIndicator
func (spi *SimpleProgressIndicator) GetName() string {
	return spi.name
}

// GetType implements ProgressIndicator
func (spi *SimpleProgressIndicator) GetType() IndicatorType {
	return spi.indicatorType
}

// IsComplete implements ProgressIndicator
func (spi *SimpleProgressIndicator) IsComplete(ctx context.Context, job *Job) (bool, error) {
	progress, err := spi.Measure(ctx, job)
	if err != nil {
		return false, err
	}
	return progress >= 1.0, nil
}

// SimpleJobTest is a basic implementation of JobTest for demonstration
type SimpleJobTest struct {
	name        string
	description string
	testFunc    func(context.Context, *Job) (*TestResult, error)
}

// NewSimpleJobTest creates a new SimpleJobTest
func NewSimpleJobTest(name, description string, testFunc func(context.Context, *Job) (*TestResult, error)) *SimpleJobTest {
	return &SimpleJobTest{
		name:        name,
		description: description,
		testFunc:    testFunc,
	}
}

// Execute implements JobTest
func (sjt *SimpleJobTest) Execute(ctx context.Context, job *Job) (*TestResult, error) {
	if sjt.testFunc == nil {
		return nil, NewJTBDError(ErrCodeInvalidTest, "test function not set", nil)
	}
	return sjt.testFunc(ctx, job)
}

// GetTestName implements JobTest
func (sjt *SimpleJobTest) GetTestName() string {
	return sjt.name
}

// GetDescription implements JobTest
func (sjt *SimpleJobTest) GetDescription() string {
	return sjt.description
}

// Validate implements JobTest
func (sjt *SimpleJobTest) Validate() error {
	if sjt.name == "" {
		return NewJTBDError(ErrCodeInvalidTest, "test name cannot be empty", nil)
	}
	if sjt.testFunc == nil {
		return NewJTBDError(ErrCodeInvalidTest, "test function cannot be nil", nil)
	}
	return nil
}

// ExampleWalmartPantryStocking demonstrates a complete JTBD definition for Walmart
func ExampleWalmartPantryStocking() (*Job, error) {
	// Build the job using the fluent builder API
	job, err := NewJobBuilder("walmart-monthly-pantry-stock", "Stock pantry for the month").
		WithDescription("Purchase and organize groceries to feed family for one month").
		WithFunctional("Get enough groceries into my house to feed my family for a month").
		WithEmotional("Feel confident that my family won't run out of essential food items").
		WithSocial("Be seen as a reliable and organized provider by my family").
		WithIndustry("retail").
		WithCompany("walmart").
		AddCircumstance(&Circumstance{
			Type:        CircumstanceTypeTemporal,
			Description: "End of month before next paycheck",
			Constraints: map[string]interface{}{
				"days_until_paycheck": 3,
				"current_pantry_level": "20%",
			},
			Triggers: []string{
				"Calendar reminder for monthly shopping",
				"Running low on staple items",
			},
			Intensity: 0.8,
		}).
		AddCircumstance(&Circumstance{
			Type:        CircumstanceTypeSituational,
			Description: "Family of 4 with dietary restrictions",
			Constraints: map[string]interface{}{
				"family_size":         4,
				"has_restrictions":    true,
				"budget_limit":        500.00,
			},
			Intensity: 0.6,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeSpeed,
			Description: "Complete shopping trip in under 90 minutes",
			Metric:      "total_shopping_time_minutes",
			Target:      90.0,
			Unit:        "minutes",
			Priority:    2,
			Direction:   "minimize",
			Threshold:   120.0,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeCost,
			Description: "Stay within monthly grocery budget",
			Metric:      "total_cost_dollars",
			Target:      450.00,
			Unit:        "dollars",
			Priority:    1,
			Direction:   "minimize",
			Threshold:   500.00,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeQuality,
			Description: "Get all items on shopping list",
			Metric:      "list_completion_percentage",
			Target:      100.0,
			Unit:        "percent",
			Priority:    1,
			Direction:   "maximize",
			Threshold:   90.0,
		}).
		WithMetadata("customer_segment", "suburban_family").
		WithMetadata("shopping_frequency", "monthly").
		Build()

	return job, err
}

// ExampleAmazonGiftFinding demonstrates a JTBD definition for Amazon gift shopping
func ExampleAmazonGiftFinding() (*Job, error) {
	job, err := NewJobBuilder("amazon-quick-gift-find", "Find perfect gift quickly").
		WithDescription("Discover and purchase an appropriate gift under time pressure").
		WithFunctional("Find and order a gift that the recipient will appreciate").
		WithEmotional("Feel confident that I've chosen a thoughtful, appropriate gift").
		WithSocial("Be seen as a thoughtful friend who remembers important occasions").
		WithIndustry("e-commerce").
		WithCompany("amazon").
		AddCircumstance(&Circumstance{
			Type:        CircumstanceTypeTemporal,
			Description: "Birthday is in 2 days",
			Constraints: map[string]interface{}{
				"days_until_event":   2,
				"needs_prime_shipping": true,
			},
			Triggers: []string{
				"Calendar notification for upcoming birthday",
				"Social media reminder",
			},
			Intensity: 0.9,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeSpeed,
			Description: "Find suitable gift in under 15 minutes",
			Metric:      "search_time_minutes",
			Target:      15.0,
			Unit:        "minutes",
			Priority:    1,
			Direction:   "minimize",
			Threshold:   30.0,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeExperience,
			Description: "Feel confident about gift choice",
			Metric:      "confidence_score",
			Target:      0.8,
			Unit:        "score_0_to_1",
			Priority:    1,
			Direction:   "maximize",
			Threshold:   0.6,
		}).
		WithMetadata("occasion", "birthday").
		WithMetadata("relationship", "friend").
		Build()

	return job, err
}

// ExampleAppleFamilyConnection demonstrates a JTBD definition for Apple device ecosystem
func ExampleAppleFamilyConnection() (*Job, error) {
	job, err := NewJobBuilder("apple-family-connection", "Stay connected with family").
		WithDescription("Maintain seamless communication and content sharing with family members").
		WithFunctional("Exchange messages, photos, and video calls with family members daily").
		WithEmotional("Feel close to family members even when physically distant").
		WithSocial("Be seen as an engaged family member who stays in touch").
		WithIndustry("technology").
		WithCompany("apple").
		AddCircumstance(&Circumstance{
			Type:        CircumstanceTypeSpatial,
			Description: "Family members live in different cities",
			Constraints: map[string]interface{}{
				"distance_miles":      500,
				"time_zone_difference": 2,
			},
			Intensity: 0.7,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeSpeed,
			Description: "Connect with family member in under 10 seconds",
			Metric:      "call_initiation_time_seconds",
			Target:      10.0,
			Unit:        "seconds",
			Priority:    2,
			Direction:   "minimize",
			Threshold:   30.0,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeQuality,
			Description: "Maintain high quality video and audio",
			Metric:      "call_quality_score",
			Target:      0.9,
			Unit:        "score_0_to_1",
			Priority:    1,
			Direction:   "maximize",
			Threshold:   0.7,
		}).
		WithMetadata("usage_pattern", "daily").
		Build()

	return job, err
}

// ExampleCVSPrescriptionManagement demonstrates a JTBD definition for CVS Health
func ExampleCVSPrescriptionManagement() (*Job, error) {
	job, err := NewJobBuilder("cvs-prescription-management", "Manage prescriptions easily").
		WithDescription("Keep track of and refill multiple prescriptions without confusion").
		WithFunctional("Refill all my prescriptions on time without errors").
		WithEmotional("Feel confident that I have the right medications when I need them").
		WithSocial("Be seen as responsible about my health by my family and doctor").
		WithIndustry("healthcare").
		WithCompany("cvs").
		AddCircumstance(&Circumstance{
			Type:        CircumstanceTypeSituational,
			Description: "Managing 5 different prescriptions with different schedules",
			Constraints: map[string]interface{}{
				"prescription_count": 5,
				"refill_frequencies": []string{"monthly", "monthly", "quarterly", "monthly", "as-needed"},
			},
			Intensity: 0.8,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeSpeed,
			Description: "Refill prescription in under 5 minutes",
			Metric:      "refill_time_minutes",
			Target:      5.0,
			Unit:        "minutes",
			Priority:    1,
			Direction:   "minimize",
			Threshold:   10.0,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeQuality,
			Description: "Zero prescription errors or mix-ups",
			Metric:      "error_rate",
			Target:      0.0,
			Unit:        "errors_per_100_refills",
			Priority:    1,
			Direction:   "minimize",
			Threshold:   1.0,
		}).
		WithMetadata("patient_type", "chronic_condition").
		Build()

	return job, err
}

// ExampleUnitedHealthCoverageUnderstanding demonstrates a JTBD definition for UnitedHealth
func ExampleUnitedHealthCoverageUnderstanding() (*Job, error) {
	job, err := NewJobBuilder("unitedhealth-coverage-understanding", "Understand medical coverage").
		WithDescription("Determine what medical procedures and costs are covered before scheduling").
		WithFunctional("Find out exact coverage and out-of-pocket costs for an upcoming procedure").
		WithEmotional("Feel confident and stress-free about the financial aspects of my healthcare").
		WithSocial("Be seen as financially responsible and informed by my family").
		WithIndustry("healthcare").
		WithCompany("unitedhealth").
		AddCircumstance(&Circumstance{
			Type:        CircumstanceTypeTemporal,
			Description: "Medical procedure scheduled in 1 week",
			Constraints: map[string]interface{}{
				"days_until_procedure": 7,
				"procedure_type":       "elective surgery",
			},
			Intensity: 0.9,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeSpeed,
			Description: "Get coverage answer in under 10 minutes",
			Metric:      "time_to_answer_minutes",
			Target:      10.0,
			Unit:        "minutes",
			Priority:    2,
			Direction:   "minimize",
			Threshold:   60.0,
		}).
		AddOutcome(&Outcome{
			Type:        OutcomeTypeQuality,
			Description: "Receive accurate cost estimate within 5% of actual",
			Metric:      "cost_estimate_accuracy_percentage",
			Target:      95.0,
			Unit:        "percent",
			Priority:    1,
			Direction:   "maximize",
			Threshold:   85.0,
		}).
		WithMetadata("coverage_type", "PPO").
		Build()

	return job, err
}

// ExampleBasicTestExecution demonstrates how to execute tests with the framework
func ExampleBasicTestExecution() error {
	// Create registry and executor
	registry := NewJobRegistry()
	executor := NewTestExecutor(registry)

	// Create and register a Walmart job
	walmartJob, err := ExampleWalmartPantryStocking()
	if err != nil {
		return err
	}

	if err := registry.RegisterJob(walmartJob); err != nil {
		return err
	}

	// Create a simple test that checks outcome achievement
	test := NewSimpleJobTest(
		"outcome_achievement_test",
		"Validates that all priority 1 outcomes meet their thresholds",
		func(ctx context.Context, job *Job) (*TestResult, error) {
			result := &TestResult{
				TestName:             "outcome_achievement_test",
				JobID:                job.ID,
				Success:              true,
				Score:                1.0,
				ProgressMeasurements: make(map[string]float64),
				OutcomeResults:       make(map[string]*OutcomeResult),
				Timestamp:            time.Now(),
				Metadata:             make(map[string]interface{}),
			}

			// In a real test, you would measure actual values
			// For this example, we'll simulate values
			for _, outcome := range job.Outcomes {
				// Simulate measurement (in real usage, this would call actual measurement logic)
				actualValue := outcome.Target * 0.95 // Assume 95% of target achieved

				outcomeResult := &OutcomeResult{
					OutcomeDescription: outcome.Description,
					MetricName:         outcome.Metric,
					ActualValue:        actualValue,
					TargetValue:        outcome.Target,
					ThresholdValue:     outcome.Threshold,
					Unit:               outcome.Unit,
					MetThreshold:       true,
					MetTarget:          false,
					PerformanceRatio:   actualValue / outcome.Target,
				}

				result.OutcomeResults[outcome.Metric] = outcomeResult

				if outcome.Priority == 1 && !outcomeResult.MetThreshold {
					result.Success = false
					result.Score = 0.0
				}
			}

			if result.Success {
				result.Message = "All priority 1 outcomes met their thresholds"
			} else {
				result.Message = "One or more priority 1 outcomes did not meet thresholds"
			}

			return result, nil
		},
	)

	// Register and execute the test
	if err := executor.RegisterTest(test); err != nil {
		return err
	}

	ctx := context.Background()
	testResult, err := executor.ExecuteTest(ctx, "outcome_achievement_test", walmartJob.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Test Result: %s\n", testResult.Message)
	fmt.Printf("Success: %v, Score: %.2f\n", testResult.Success, testResult.Score)

	return nil
}

// ExampleProgressIndicatorUsage demonstrates creating and using progress indicators
func ExampleProgressIndicatorUsage() (*Job, error) {
	// Create a job
	job, err := NewJobBuilder("example-job", "Example Job").
		WithDescription("Demonstrates progress indicator usage").
		WithFunctional("Complete a task").
		WithEmotional("Feel accomplished").
		WithSocial("Be recognized for completion").
		Build()

	if err != nil {
		return nil, err
	}

	// Create a leading indicator (predicts success)
	leadingIndicator := NewSimpleProgressIndicator(
		"task_started",
		IndicatorTypeLeading,
		func(ctx context.Context, j *Job) (float64, error) {
			// In real usage, check if task has started
			// Return 1.0 if started, 0.0 if not
			return 1.0, nil
		},
	)

	// Create a concurrent indicator (tracks during execution)
	concurrentIndicator := NewSimpleProgressIndicator(
		"steps_completed",
		IndicatorTypeConcurrent,
		func(ctx context.Context, j *Job) (float64, error) {
			// In real usage, calculate percentage of steps completed
			// For example: completed_steps / total_steps
			return 0.65, nil
		},
	)

	// Create a lagging indicator (measures final outcome)
	laggingIndicator := NewSimpleProgressIndicator(
		"quality_check",
		IndicatorTypeLagging,
		func(ctx context.Context, j *Job) (float64, error) {
			// In real usage, perform quality check after completion
			// Return quality score
			return 0.92, nil
		},
	)

	// Add indicators to job
	job.Indicators = []ProgressIndicator{
		leadingIndicator,
		concurrentIndicator,
		laggingIndicator,
	}

	return job, nil
}
