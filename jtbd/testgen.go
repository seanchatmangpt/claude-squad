package jtbd

import (
	"fmt"
	"strings"
	"time"
)

// TestJobSpec represents job specifications for test case generation
type TestJobSpec struct {
	Name        string
	Description string
	Category    string
	Steps       []string
	Priority    string
	Functional  string
	Emotional   string
	Social      string
}

// TestCircumstanceSpec represents circumstance specifications for test generation
type TestCircumstanceSpec struct {
	Location    string
	TimeOfDay   string
	Season      string
	Urgency     string
	Environment string
	Triggers    []string
	Intensity   float64
}

// TestOutcomeSpec represents outcome specifications for test generation
type TestOutcomeSpec struct {
	Success     bool
	Description string
	Metrics     map[string]interface{}
	SideEffects []string
	Type        OutcomeType
	Target      float64
	Unit        string
}

// Constraint represents limitations or requirements for test cases
type Constraint struct {
	Type        string
	Description string
	Value       interface{}
	Hard        bool
}

// TestCase represents a complete JTBD test scenario
type TestCase struct {
	ID               string
	Industry         string
	JobSpec          TestJobSpec
	CircumstanceSpec TestCircumstanceSpec
	OutcomeSpec      TestOutcomeSpec
	Constraints      []Constraint
	CompetingJobs    []TestJobSpec
	TradeOffs        []string
	Variations       []string
	IsEdgeCase       bool
	IsHappyPath      bool
	MultiStep        bool
	StepSequence     []string
}

// ToJob converts a TestCase into a framework Job
func (tc *TestCase) ToJob() *Job {
	job := &Job{
		ID:          tc.ID,
		Name:        tc.JobSpec.Name,
		Description: tc.JobSpec.Description,
		Functional:  tc.JobSpec.Functional,
		Emotional:   tc.JobSpec.Emotional,
		Social:      tc.JobSpec.Social,
		Industry:    tc.Industry,
		Metadata:    make(map[string]interface{}),
	}

	// Add circumstance
	if tc.CircumstanceSpec.Triggers != nil || tc.CircumstanceSpec.Urgency != "" {
		circ := &Circumstance{
			Type:        CircumstanceTypeSituational,
			Description: fmt.Sprintf("%s at %s", tc.CircumstanceSpec.Urgency, tc.CircumstanceSpec.TimeOfDay),
			Constraints: make(map[string]interface{}),
			Triggers:    tc.CircumstanceSpec.Triggers,
			Intensity:   tc.CircumstanceSpec.Intensity,
		}
		if tc.CircumstanceSpec.Location != "" {
			circ.Constraints["location"] = tc.CircumstanceSpec.Location
		}
		job.Circumstances = append(job.Circumstances, circ)
	}

	// Add outcome
	if tc.OutcomeSpec.Description != "" {
		outcome := &Outcome{
			Type:        tc.OutcomeSpec.Type,
			Description: tc.OutcomeSpec.Description,
			Target:      tc.OutcomeSpec.Target,
			Unit:        tc.OutcomeSpec.Unit,
			Metadata:    tc.OutcomeSpec.Metrics,
		}
		job.Outcomes = append(job.Outcomes, outcome)
	}

	// Add metadata
	job.Metadata["test_case_id"] = tc.ID
	job.Metadata["category"] = tc.JobSpec.Category
	job.Metadata["steps"] = tc.JobSpec.Steps

	return job
}

// TestCaseGenerator generates comprehensive JTBD test cases
type TestCaseGenerator struct {
	industryPatterns map[string]*IndustryPattern
	testCaseCounter  int
}

// IndustryPattern defines patterns for specific industries
type IndustryPattern struct {
	Name       string
	Jobs       []JobTemplate
	Outcomes   []OutcomeTemplate
}

// JobTemplate is a template for generating jobs
type JobTemplate struct {
	Name        string
	Description string
	Category    string
	Steps       []string
	Priority    string
	Functional  string
	Emotional   string
	Social      string
}

// OutcomeTemplate is a template for generating outcomes
type OutcomeTemplate struct {
	Success     bool
	Description string
	Type        OutcomeType
	Target      float64
	Unit        string
}

// NewTestCaseGenerator creates a new test case generator
func NewTestCaseGenerator() *TestCaseGenerator {
	gen := &TestCaseGenerator{
		industryPatterns: make(map[string]*IndustryPattern),
		testCaseCounter:  0,
	}
	gen.initializePatterns()
	return gen
}

// initializePatterns sets up all Fortune 5 industry patterns
func (g *TestCaseGenerator) initializePatterns() {
	g.industryPatterns["retail"] = g.createRetailPattern()
	g.industryPatterns["ecommerce"] = g.createEcommercePattern()
	g.industryPatterns["technology"] = g.createTechnologyPattern()
	g.industryPatterns["healthcare"] = g.createHealthcarePattern()
	g.industryPatterns["insurance"] = g.createInsurancePattern()
}

// createRetailPattern creates Walmart-style retail patterns
func (g *TestCaseGenerator) createRetailPattern() *IndustryPattern {
	return &IndustryPattern{
		Name: "Retail (Walmart)",
		Jobs: []JobTemplate{
			{
				Name:        "Weekly Grocery Shopping",
				Description: "Purchase groceries for the week",
				Category:    "procurement",
				Functional:  "Get enough groceries to feed family for a week",
				Emotional:   "Feel confident family has enough food",
				Social:      "Be seen as responsible provider",
				Steps:       []string{"Create list", "Shop", "Checkout", "Transport home"},
				Priority:    "high",
			},
		},
		Outcomes: []OutcomeTemplate{
			{
				Success:     true,
				Description: "Shopping completed within budget and time",
				Type:        OutcomeTypeSpeed,
				Target:      35.0,
				Unit:        "minutes",
			},
		},
	}
}

// createEcommercePattern creates Amazon-style patterns
func (g *TestCaseGenerator) createEcommercePattern() *IndustryPattern {
	return &IndustryPattern{
		Name: "E-commerce (Amazon)",
		Jobs: []JobTemplate{
			{
				Name:        "Gift Shopping",
				Description: "Find and purchase gift",
				Category:    "gifting",
				Functional:  "Find perfect gift within budget",
				Emotional:   "Feel confident recipient will love it",
				Social:      "Be seen as thoughtful",
				Steps:       []string{"Search", "Compare", "Select", "Purchase"},
				Priority:    "high",
			},
		},
		Outcomes: []OutcomeTemplate{
			{
				Success:     true,
				Description: "Gift found and delivered on time",
				Type:        OutcomeTypeQuality,
				Target:      4.5,
				Unit:        "rating",
			},
		},
	}
}

// createTechnologyPattern creates Apple-style patterns
func (g *TestCaseGenerator) createTechnologyPattern() *IndustryPattern {
	return &IndustryPattern{
		Name: "Technology (Apple)",
		Jobs: []JobTemplate{
			{
				Name:        "Device Setup for Non-Technical User",
				Description: "Set up new device for someone with limited tech skills",
				Category:    "setup assistance",
				Functional:  "Get device working with data transferred",
				Emotional:   "Feel confident they can use it",
				Social:      "Be seen as helpful by family",
				Steps:       []string{"Unbox", "Transfer data", "Configure", "Teach basics"},
				Priority:    "high",
			},
		},
		Outcomes: []OutcomeTemplate{
			{
				Success:     true,
				Description: "Device setup complete with user confident",
				Type:        OutcomeTypeExperience,
				Target:      4.0,
				Unit:        "confidence_rating",
			},
		},
	}
}

// createHealthcarePattern creates CVS Health-style patterns
func (g *TestCaseGenerator) createHealthcarePattern() *IndustryPattern {
	return &IndustryPattern{
		Name: "Healthcare/Pharmacy (CVS Health)",
		Jobs: []JobTemplate{
			{
				Name:        "Prescription Refill",
				Description: "Refill recurring prescription",
				Category:    "medication management",
				Functional:  "Get medication refilled without running out",
				Emotional:   "Feel secure about health management",
				Social:      "Maintain independence in managing health",
				Steps:       []string{"Check inventory", "Request refill", "Pick up"},
				Priority:    "high",
			},
		},
		Outcomes: []OutcomeTemplate{
			{
				Success:     true,
				Description: "Prescription refilled successfully",
				Type:        OutcomeTypeSpeed,
				Target:      10.0,
				Unit:        "minutes",
			},
		},
	}
}

// createInsurancePattern creates UnitedHealth Group-style patterns
func (g *TestCaseGenerator) createInsurancePattern() *IndustryPattern {
	return &IndustryPattern{
		Name: "Insurance (UnitedHealth Group)",
		Jobs: []JobTemplate{
			{
				Name:        "Find In-Network Provider",
				Description: "Locate provider covered by insurance",
				Category:    "provider search",
				Functional:  "Find qualified provider who accepts insurance",
				Emotional:   "Feel confident about costs",
				Social:      "Make responsible healthcare decisions",
				Steps:       []string{"Access directory", "Filter", "Verify", "Schedule"},
				Priority:    "high",
			},
		},
		Outcomes: []OutcomeTemplate{
			{
				Success:     true,
				Description: "Found provider and scheduled appointment",
				Type:        OutcomeTypeSpeed,
				Target:      15.0,
				Unit:        "minutes",
			},
		},
	}
}

// TestGenerationOptions configures test case generation
type TestGenerationOptions struct {
	IncludeHappyPath    bool
	IncludeEdgeCases    bool
	IncludeFailures     bool
	IncludeMultiStep    bool
	IncludeCompeting    bool
	CombinatorialLevel  int
	MaxCasesPerCategory int
}

// GenerateTestCases generates test cases for a specific industry
func (g *TestCaseGenerator) GenerateTestCases(industry string, options TestGenerationOptions) []TestCase {
	pattern, exists := g.industryPatterns[strings.ToLower(industry)]
	if !exists {
		return nil
	}

	var testCases []TestCase

	if options.IncludeHappyPath {
		testCases = append(testCases, g.generateHappyPathCases(pattern)...)
	}

	if options.IncludeEdgeCases {
		testCases = append(testCases, g.generateEdgeCases(pattern)...)
	}

	if options.IncludeFailures {
		testCases = append(testCases, g.generateFailureCases(pattern)...)
	}

	if options.IncludeMultiStep {
		testCases = append(testCases, g.generateMultiStepCases(pattern)...)
	}

	if options.IncludeCompeting {
		testCases = append(testCases, g.generateCompetingJobsCases(pattern)...)
	}

	if options.CombinatorialLevel > 0 {
		testCases = g.explodeCombinations(testCases, options.CombinatorialLevel)
	}

	return testCases
}

// generateHappyPathCases generates standard success scenarios
func (g *TestCaseGenerator) generateHappyPathCases(pattern *IndustryPattern) []TestCase {
	var cases []TestCase

	for _, jobTemplate := range pattern.Jobs {
		tc := TestCase{
			ID:          g.nextID(),
			Industry:    pattern.Name,
			IsHappyPath: true,
			JobSpec:     TestJobSpec{
				Name:        jobTemplate.Name,
				Description: jobTemplate.Description,
				Category:    jobTemplate.Category,
				Steps:       jobTemplate.Steps,
				Priority:    jobTemplate.Priority,
				Functional:  jobTemplate.Functional,
				Emotional:   jobTemplate.Emotional,
				Social:      jobTemplate.Social,
			},
			CircumstanceSpec: TestCircumstanceSpec{
				Urgency:   "normal",
				Intensity: 0.5,
			},
		}

		if len(pattern.Outcomes) > 0 {
			outcome := pattern.Outcomes[0]
			tc.OutcomeSpec = TestOutcomeSpec{
				Success:     outcome.Success,
				Description: outcome.Description,
				Type:        outcome.Type,
				Target:      outcome.Target,
				Unit:        outcome.Unit,
			}
		}

		cases = append(cases, tc)
	}

	return cases
}

// generateEdgeCases generates edge case scenarios
func (g *TestCaseGenerator) generateEdgeCases(pattern *IndustryPattern) []TestCase {
	var cases []TestCase

	for _, jobTemplate := range pattern.Jobs {
		tc := TestCase{
			ID:          g.nextID(),
			Industry:    pattern.Name,
			IsEdgeCase:  true,
			JobSpec:     TestJobSpec{
				Name:        jobTemplate.Name,
				Description: jobTemplate.Description,
				Category:    jobTemplate.Category,
				Steps:       jobTemplate.Steps,
				Priority:    "high",
				Functional:  jobTemplate.Functional,
				Emotional:   jobTemplate.Emotional,
				Social:      jobTemplate.Social,
			},
			CircumstanceSpec: TestCircumstanceSpec{
				Urgency:   "immediate",
				Intensity: 0.9,
			},
			Constraints: []Constraint{
				{
					Type:        "time",
					Description: "Must complete within 15 minutes",
					Value:       15,
					Hard:        true,
				},
			},
			TradeOffs: []string{"willing to pay premium"},
		}

		cases = append(cases, tc)
	}

	return cases
}

// generateFailureCases generates failure scenarios
func (g *TestCaseGenerator) generateFailureCases(pattern *IndustryPattern) []TestCase {
	var cases []TestCase

	for _, jobTemplate := range pattern.Jobs {
		tc := TestCase{
			ID:          g.nextID(),
			Industry:    pattern.Name,
			JobSpec:     TestJobSpec{
				Name:        jobTemplate.Name,
				Description: jobTemplate.Description,
				Category:    jobTemplate.Category,
				Steps:       jobTemplate.Steps,
				Priority:    jobTemplate.Priority,
				Functional:  jobTemplate.Functional,
				Emotional:   jobTemplate.Emotional,
				Social:      jobTemplate.Social,
			},
			OutcomeSpec: TestOutcomeSpec{
				Success:     false,
				Description: "Failed to complete job",
			},
		}

		cases = append(cases, tc)
	}

	return cases
}

// generateMultiStepCases generates complex multi-step workflows
func (g *TestCaseGenerator) generateMultiStepCases(pattern *IndustryPattern) []TestCase {
	var cases []TestCase

	if len(pattern.Jobs) >= 2 {
		tc := TestCase{
			ID:          g.nextID(),
			Industry:    pattern.Name,
			IsHappyPath: true,
			MultiStep:   true,
			JobSpec:     TestJobSpec{
				Name:        "Multi-Job Workflow",
				Description: "Complex workflow with multiple jobs",
				Category:    "workflow",
				Priority:    "high",
				Functional:  "Complete multiple related tasks efficiently",
				Emotional:   "Feel accomplished",
				Social:      "Be seen as efficient",
			},
			StepSequence: []string{pattern.Jobs[0].Name, pattern.Jobs[1].Name},
		}

		cases = append(cases, tc)
	}

	return cases
}

// generateCompetingJobsCases generates scenarios with competing priorities
func (g *TestCaseGenerator) generateCompetingJobsCases(pattern *IndustryPattern) []TestCase {
	var cases []TestCase

	if len(pattern.Jobs) >= 2 {
		tc := TestCase{
			ID:          g.nextID(),
			Industry:    pattern.Name,
			IsEdgeCase:  true,
			MultiStep:   true,
			JobSpec:     TestJobSpec{
				Name:        pattern.Jobs[0].Name,
				Description: pattern.Jobs[0].Description,
				Category:    pattern.Jobs[0].Category,
				Steps:       pattern.Jobs[0].Steps,
				Priority:    "high",
				Functional:  pattern.Jobs[0].Functional,
				Emotional:   pattern.Jobs[0].Emotional,
				Social:      pattern.Jobs[0].Social,
			},
			CompetingJobs: []TestJobSpec{
				{
					Name:        pattern.Jobs[1].Name,
					Description: pattern.Jobs[1].Description,
					Category:    pattern.Jobs[1].Category,
					Priority:    "high",
					Functional:  pattern.Jobs[1].Functional,
					Emotional:   pattern.Jobs[1].Emotional,
					Social:      pattern.Jobs[1].Social,
				},
			},
			TradeOffs: []string{"prioritize primary job"},
		}

		cases = append(cases, tc)
	}

	return cases
}

// explodeCombinations creates combinatorial variations
func (g *TestCaseGenerator) explodeCombinations(baseCases []TestCase, level int) []TestCase {
	if level == 0 {
		return baseCases
	}

	var exploded []TestCase
	exploded = append(exploded, baseCases...)

	for _, baseCase := range baseCases {
		for i := 0; i < level * 2; i++ {
			variant := baseCase
			variant.ID = g.nextID()
			variant.Variations = append(variant.Variations, fmt.Sprintf("variation-%d", i))
			exploded = append(exploded, variant)
		}
	}

	return exploded
}

func (g *TestCaseGenerator) nextID() string {
	g.testCaseCounter++
	timestamp := time.Now().Format("20060102")
	return fmt.Sprintf("TC-%s-%04d", timestamp, g.testCaseCounter)
}

// GetAllIndustries returns all supported industries
func (g *TestCaseGenerator) GetAllIndustries() []string {
	return []string{"retail", "ecommerce", "technology", "healthcare", "insurance"}
}

// GetIndustryPattern returns the pattern for a specific industry
func (g *TestCaseGenerator) GetIndustryPattern(industry string) *IndustryPattern {
	return g.industryPatterns[strings.ToLower(industry)]
}

// GenerateAllTestCases generates test cases for all industries
func (g *TestCaseGenerator) GenerateAllTestCases(options TestGenerationOptions) map[string][]TestCase {
	allCases := make(map[string][]TestCase)

	for industry := range g.industryPatterns {
		allCases[industry] = g.GenerateTestCases(industry, options)
	}

	return allCases
}
