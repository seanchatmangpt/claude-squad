package jtbd

import (
	"testing"
)

func TestNewTestCaseGenerator(t *testing.T) {
	gen := NewTestCaseGenerator()

	if gen == nil {
		t.Fatal("Expected generator to be created, got nil")
	}

	if len(gen.industryPatterns) != 5 {
		t.Errorf("Expected 5 industry patterns, got %d", len(gen.industryPatterns))
	}
}

func TestGetAllIndustries(t *testing.T) {
	gen := NewTestCaseGenerator()
	industries := gen.GetAllIndustries()

	if len(industries) != 5 {
		t.Errorf("Expected 5 industries, got %d", len(industries))
	}
}

func TestGenerateRetailTestCases(t *testing.T) {
	gen := NewTestCaseGenerator()

	options := TestGenerationOptions{
		IncludeHappyPath: true,
		IncludeEdgeCases: true,
		IncludeFailures:  true,
	}

	cases := gen.GenerateTestCases("retail", options)

	if len(cases) == 0 {
		t.Fatal("Expected test cases to be generated, got none")
	}

	hasHappyPath := false
	hasEdgeCase := false

	for _, tc := range cases {
		if tc.IsHappyPath {
			hasHappyPath = true
		}
		if tc.IsEdgeCase {
			hasEdgeCase = true
		}

		if tc.ID == "" {
			t.Error("Test case ID should not be empty")
		}
		if tc.Industry == "" {
			t.Error("Test case industry should not be empty")
		}
		if tc.JobSpec.Name == "" {
			t.Error("Job name should not be empty")
		}
		if tc.JobSpec.Functional == "" {
			t.Error("Functional dimension should be set")
		}
	}

	if !hasHappyPath {
		t.Error("Expected at least one happy path test case")
	}
	if !hasEdgeCase {
		t.Error("Expected at least one edge case")
	}

	t.Logf("Generated %d retail test cases", len(cases))
}

func TestGenerateAllIndustries(t *testing.T) {
	gen := NewTestCaseGenerator()

	options := TestGenerationOptions{
		IncludeHappyPath: true,
		IncludeEdgeCases: true,
	}

	allCases := gen.GenerateAllTestCases(options)

	if len(allCases) != 5 {
		t.Errorf("Expected 5 industries, got %d", len(allCases))
	}

	totalCases := 0
	for industry, cases := range allCases {
		if len(cases) == 0 {
			t.Errorf("Expected test cases for %s", industry)
		}
		totalCases += len(cases)
		t.Logf("%s: %d cases", industry, len(cases))
	}

	if totalCases == 0 {
		t.Fatal("Expected total cases > 0")
	}

	t.Logf("Total test cases generated: %d", totalCases)
}

func TestTestCaseToJob(t *testing.T) {
	gen := NewTestCaseGenerator()

	options := TestGenerationOptions{
		IncludeHappyPath: true,
	}

	cases := gen.GenerateTestCases("retail", options)

	if len(cases) == 0 {
		t.Fatal("Expected test cases")
	}

	testCase := cases[0]
	job := testCase.ToJob()

	if job == nil {
		t.Fatal("Expected job conversion to succeed")
	}

	if job.ID != testCase.ID {
		t.Errorf("Expected job ID %s, got %s", testCase.ID, job.ID)
	}

	if job.Functional != testCase.JobSpec.Functional {
		t.Error("Functional dimension should match")
	}

	if job.Emotional != testCase.JobSpec.Emotional {
		t.Error("Emotional dimension should match")
	}

	if job.Social != testCase.JobSpec.Social {
		t.Error("Social dimension should match")
	}

	t.Logf("Successfully converted TestCase to Job")
}

func TestCombinatorialExplosion(t *testing.T) {
	gen := NewTestCaseGenerator()

	options1 := TestGenerationOptions{
		IncludeHappyPath:   true,
		CombinatorialLevel: 0,
	}
	cases1 := gen.GenerateTestCases("ecommerce", options1)

	options2 := TestGenerationOptions{
		IncludeHappyPath:   true,
		CombinatorialLevel: 2,
	}
	cases2 := gen.GenerateTestCases("ecommerce", options2)

	if len(cases2) <= len(cases1) {
		t.Errorf("Expected more cases with combinatorial explosion: %d vs %d", len(cases2), len(cases1))
	}

	t.Logf("Combinatorial level 0: %d cases, level 2: %d cases", len(cases1), len(cases2))
}
