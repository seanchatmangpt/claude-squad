package jtbd

import "fmt"

// ExampleDataFactory demonstrates the usage of the DataFactory
func ExampleDataFactory() {
	// Create the data factory
	df := NewDataFactory()

	// Get statistics
	stats := df.GetStatistics()
	fmt.Printf("Data Factory Statistics:\n")
	fmt.Printf("Total Personas: %d\n", stats["total_personas"])
	fmt.Printf("Products by Company: %v\n", stats["products_by_company"])

	// Example 1: Walmart Grocery Scenario
	fmt.Println("\n=== Walmart Grocery Scenario ===")
	walmartScenario := df.GetWalmartGroceryScenario("sarah_budget")
	persona := walmartScenario["persona"].(*Persona)
	products := walmartScenario["products"].([]*Product)
	context := walmartScenario["context"].(*Context)

	fmt.Printf("Persona: %s (%d years old, $%d income)\n", persona.Name, persona.Age, persona.Income)
	fmt.Printf("Segment: %s | Price Sensitivity: %s\n", persona.Segment, persona.PriceSensitivity)
	fmt.Printf("Time Context: %s\n", context.TimeContext)
	fmt.Printf("Products selected: %d\n", len(products))
	for _, p := range products {
		fmt.Printf("  - %s ($%.2f)\n", p.Name, p.Price)
	}

	// Example 2: Amazon Prime Scenario
	fmt.Println("\n=== Amazon Prime Scenario ===")
	amazonScenario := df.GetAmazonPrimeScenario("tyler_techsavvy")
	persona = amazonScenario["persona"].(*Persona)
	products = amazonScenario["products"].([]*Product)

	fmt.Printf("Persona: %s | Segment: %s\n", persona.Name, persona.Segment)
	fmt.Printf("Tech Savviness: %s\n", persona.TechSavviness)
	for _, p := range products {
		fmt.Printf("  - %s ($%.2f, Rating: %.1f)\n", p.Name, p.Price, p.Rating)
	}

	// Example 3: Apple Ecosystem Scenario
	fmt.Println("\n=== Apple Ecosystem Scenario ===")
	appleScenario := df.GetAppleEcosystemScenario("patricia_premium")
	persona = appleScenario["persona"].(*Persona)
	products = appleScenario["products"].([]*Product)
	context = appleScenario["context"].(*Context)

	fmt.Printf("Persona: %s | Income: $%d | Family: %d\n", persona.Name, persona.Income, persona.FamilySize)
	fmt.Printf("Event: %s (Urgency: %s)\n", context.EventContext.Type, context.EventContext.Urgency)
	fmt.Printf("Budget: $%.2f\n", context.Constraints.Budget)
	for _, p := range products {
		fmt.Printf("  - %s ($%.2f)\n", p.Name, p.Price)
	}

	// Example 4: CVS Pharmacy Scenario
	fmt.Println("\n=== CVS Pharmacy Scenario ===")
	cvsScenario := df.GetCVSPharmacyScenario("edward_elderly")
	persona = cvsScenario["persona"].(*Persona)
	products = cvsScenario["products"].([]*Product)

	fmt.Printf("Persona: %s | Age: %d | Segment: %s\n", persona.Name, persona.Age, persona.Segment)
	for _, p := range products {
		fmt.Printf("  - %s ($%.2f, Category: %s)\n", p.Name, p.Price, p.Category)
	}

	// Example 5: UnitedHealth Enrollment Scenario
	fmt.Println("\n=== UnitedHealth Enrollment Scenario ===")
	uhScenario := df.GetUnitedHealthEnrollmentScenario("fatima_family")
	persona = uhScenario["persona"].(*Persona)
	products = uhScenario["products"].([]*Product)

	fmt.Printf("Persona: %s | Family Size: %d\n", persona.Name, persona.FamilySize)
	fmt.Printf("Monthly Budget: $%.2f\n", uhScenario["context"].(*Context).Constraints.Budget)
	for _, p := range products {
		fmt.Printf("  - %s ($%.2f/month, Rating: %.1f)\n", p.Name, p.Price, p.Rating)
	}

	// Example 6: Using Builder Pattern
	fmt.Println("\n=== Custom Persona Builder ===")
	customPersona := NewPersonaBuilder("custom_001", "Alex Developer").
		WithAge(32).
		WithIncome(120000).
		WithFamilySize(1).
		WithLocation(Urban).
		WithSegment(TechSavvy).
		WithTechSavviness(TechExpert).
		WithPriceSensitivity(Medium).
		WithPreference("remote_work", true).
		WithPreference("cloud_services", true).
		Build()

	fmt.Printf("Created: %s | Age: %d | Income: $%d\n", customPersona.Name, customPersona.Age, customPersona.Income)
	fmt.Printf("Preferences: %v\n", customPersona.Preferences)

	// Example 7: Generate Random Transaction
	fmt.Println("\n=== Random Transaction Generation ===")
	transaction := df.GenerateRandomTransaction("sarah_budget", Walmart, 5)
	if transaction != nil {
		fmt.Printf("Transaction ID: %s\n", transaction.ID)
		fmt.Printf("Persona: %s\n", transaction.PersonaID)
		fmt.Printf("Channel: %s\n", transaction.Channel)
		fmt.Printf("Total Amount: $%.2f\n", transaction.TotalAmount)
		fmt.Printf("Products purchased: %d\n", len(transaction.Products))
		for _, purchase := range transaction.Products {
			fmt.Printf("  - %s x%d = $%.2f\n", purchase.Product.Name, purchase.Quantity, purchase.Price)
		}
	}

	// Example 8: Weekly Grocery List
	fmt.Println("\n=== Weekly Grocery List ===")
	groceryList := df.GenerateWeeklyGroceryList("fatima_family")
	if groceryList != nil {
		fmt.Printf("Family size: %d\n", df.GetPersona("fatima_family").FamilySize)
		fmt.Printf("Total: $%.2f\n", groceryList.TotalAmount)
		for _, purchase := range groceryList.Products {
			fmt.Printf("  - %s x%d\n", purchase.Product.Name, purchase.Quantity)
		}
	}

	// Example 9: Query Products by Category
	fmt.Println("\n=== Products by Category ===")
	prescriptions := df.GetProductsByCategory(CVS, "Prescription")
	fmt.Printf("CVS Prescriptions: %d items\n", len(prescriptions))
	for _, p := range prescriptions {
		fmt.Printf("  - %s ($%.2f)\n", p.Name, p.Price)
	}

	smartphones := df.GetProductsByCategory(Apple, "Smartphone")
	fmt.Printf("\nApple Smartphones: %d items\n", len(smartphones))
	for _, p := range smartphones {
		fmt.Printf("  - %s ($%.2f)\n", p.Name, p.Price)
	}

	// Example 10: Clone Persona
	fmt.Println("\n=== Clone Persona ===")
	original := df.GetPersona("tyler_techsavvy")
	clone := df.ClonePersona("tyler_techsavvy")
	fmt.Printf("Original: %s (ID: %s)\n", original.Name, original.ID)
	fmt.Printf("Clone: %s (ID: %s)\n", clone.Name, clone.ID)
	fmt.Printf("Same attributes: Age=%d, Income=$%d\n", clone.Age, clone.Income)
}
