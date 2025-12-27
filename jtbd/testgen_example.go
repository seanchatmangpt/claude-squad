//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"claude-squad/jtbd"
)

func main() {
	// Create the test case generator
	gen := jtbd.NewTestCaseGenerator()

	fmt.Println("=== JTBD Test Case Generator - Fortune 5 Examples ===\n")

	// Example 1: Walmart - Weekly grocery shopping
	fmt.Println("Example 1: Walmart - Weekly Grocery Shopping")
	fmt.Println("---")
	options := jtbd.TestGenerationOptions{
		IncludeHappyPath: true,
		IncludeEdgeCases: true,
	}
	walmartCases := gen.GenerateTestCases("retail", options)
	for _, tc := range walmartCases {
		fmt.Printf("Test Case: %s\n", tc.ID)
		fmt.Printf("  Job: %s\n", tc.JobSpec.Name)
		fmt.Printf("  Category: %s\n", tc.JobSpec.Category)
		fmt.Printf("  Functional: %s\n", tc.JobSpec.Functional)
		fmt.Printf("  Emotional: %s\n", tc.JobSpec.Emotional)
		fmt.Printf("  Social: %s\n", tc.JobSpec.Social)
		fmt.Printf("  Happy Path: %v, Edge Case: %v\n", tc.IsHappyPath, tc.IsEdgeCase)
		
		// Convert to framework Job
		job := tc.ToJob()
		fmt.Printf("  Converted to Job ID: %s\n", job.ID)
		fmt.Println()
	}

	// Example 2: Generate all test cases
	fmt.Println("\nExample 2: Generate All Industries")
	fmt.Println("---")
	allOptions := jtbd.TestGenerationOptions{
		IncludeHappyPath:   true,
		IncludeEdgeCases:   true,
		IncludeFailures:    true,
		IncludeMultiStep:   true,
		IncludeCompeting:   true,
		CombinatorialLevel: 1,
	}
	allCases := gen.GenerateAllTestCases(allOptions)
	
	totalCases := 0
	for industry, cases := range allCases {
		fmt.Printf("%s: %d test cases\n", industry, len(cases))
		totalCases += len(cases)
	}
	fmt.Printf("Total: %d test cases\n", totalCases)

	// Example 3: JSON export
	fmt.Println("\nExample 3: Export specific case to JSON")
	fmt.Println("---")
	if len(walmartCases) > 0 {
		jsonData, err := json.MarshalIndent(walmartCases[0], "", "  ")
		if err == nil {
			fmt.Println(string(jsonData))
		}
	}

	// Example 4: Industry-specific patterns
	fmt.Println("\nExample 4: Industry Patterns")
	fmt.Println("---")
	for _, industry := range gen.GetAllIndustries() {
		pattern := gen.GetIndustryPattern(industry)
		if pattern != nil {
			fmt.Printf("%s:\n", pattern.Name)
			for _, job := range pattern.Jobs {
				fmt.Printf("  - %s (%s)\n", job.Name, job.Category)
			}
		}
	}
}
