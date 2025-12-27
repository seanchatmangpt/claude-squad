// Package jtbd provides a comprehensive test data factory for Fortune 5 Jobs-to-be-Done scenarios.
package jtbd

import (
	"fmt"
	"math/rand"
	"time"
)

// Core types
type Persona struct {
	ID               string
	Name             string
	Age              int
	Income           int
	FamilySize       int
	Location         LocationType
	Segment          CustomerSegment
	TechSavviness    TechLevel
	PriceSensitivity PriceSensitivity
	Preferences      map[string]interface{}
	Behaviors        []Behavior
}

type Product struct {
	ID           string
	Name         string
	Category     string
	Brand        string
	Price        float64
	Company      Fortune5Company
	Attributes   map[string]interface{}
	Availability bool
	Rating       float64
}

type Transaction struct {
	ID          string
	PersonaID   string
	Products    []ProductPurchase
	TotalAmount float64
	Timestamp   time.Time
	Channel     Channel
	Context     *Context
}

type ProductPurchase struct {
	Product  *Product
	Quantity int
	Price    float64
}

type Context struct {
	TimeContext     TimeContext
	LocationContext LocationContext
	WeatherContext  WeatherContext
	EventContext    EventContext
	Constraints     Constraints
}

type Constraints struct {
	Budget       float64
	TimeLimit    time.Duration
	Availability map[string]bool
	Requirements []string
}

type Behavior struct {
	Type        BehaviorType
	Frequency   Frequency
	Triggers    []string
	Preferences map[string]interface{}
}

type Fortune5Company string

const (
	Walmart      Fortune5Company = "Walmart"
	Amazon       Fortune5Company = "Amazon"
	Apple        Fortune5Company = "Apple"
	CVS          Fortune5Company = "CVS"
	UnitedHealth Fortune5Company = "UnitedHealth"
)

type CustomerSegment string

const (
	BudgetConscious  CustomerSegment = "BudgetConscious"
	PremiumSeeker    CustomerSegment = "PremiumSeeker"
	Elderly          CustomerSegment = "Elderly"
	TechSavvy        CustomerSegment = "TechSavvy"
	Family           CustomerSegment = "Family"
	HealthFocused    CustomerSegment = "HealthFocused"
	ConvenienceFirst CustomerSegment = "ConvenienceFirst"
)

type LocationType string

const (
	Urban    LocationType = "Urban"
	Suburban LocationType = "Suburban"
	Rural    LocationType = "Rural"
)

type TechLevel string

const (
	TechNovice       TechLevel = "Novice"
	TechIntermediate TechLevel = "Intermediate"
	TechAdvanced     TechLevel = "Advanced"
	TechExpert       TechLevel = "Expert"
)

type PriceSensitivity string

const (
	VeryLow  PriceSensitivity = "VeryLow"
	Low      PriceSensitivity = "Low"
	Medium   PriceSensitivity = "Medium"
	High     PriceSensitivity = "High"
	VeryHigh PriceSensitivity = "VeryHigh"
)

type TimeContext string

const (
	RushHour      TimeContext = "RushHour"
	Weekend       TimeContext = "Weekend"
	HolidaySeason TimeContext = "HolidaySeason"
	Emergency     TimeContext = "Emergency"
	Routine       TimeContext = "Routine"
	LateNight     TimeContext = "LateNight"
)

type Channel string

const (
	InStore   Channel = "InStore"
	Online    Channel = "Online"
	MobileApp Channel = "MobileApp"
	Phone     Channel = "Phone"
	Curbside  Channel = "Curbside"
)

type BehaviorType string

const (
	WeeklyGrocery     BehaviorType = "WeeklyGrocery"
	ImpulsePurchase   BehaviorType = "ImpulsePurchase"
	ResearchIntensive BehaviorType = "ResearchIntensive"
	SubscriptionUser  BehaviorType = "SubscriptionUser"
	DealHunter        BehaviorType = "DealHunter"
	BrandLoyal        BehaviorType = "BrandLoyal"
)

type Frequency string

const (
	Daily   Frequency = "Daily"
	Weekly  Frequency = "Weekly"
	Monthly Frequency = "Monthly"
	Yearly  Frequency = "Yearly"
	Rare    Frequency = "Rare"
)

type LocationContext struct {
	Type     LocationType
	Distance float64
	Traffic  string
	Parking  string
}

type WeatherContext struct {
	Condition   string
	Temperature int
	Season      string
}

type EventContext struct {
	Type    string
	Urgency string
	Impact  string
}

type DataFactory struct {
	personas map[string]*Persona
	products map[Fortune5Company]map[string]*Product
	rand     *rand.Rand
}

func NewDataFactory() *DataFactory {
	df := &DataFactory{
		personas: make(map[string]*Persona),
		products: make(map[Fortune5Company]map[string]*Product),
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	df.initializePersonas()
	df.initializeProducts()
	return df
}

func (df *DataFactory) initializePersonas() {
	df.personas["sarah_budget"] = &Persona{
		ID: "sarah_budget", Name: "Sarah Martinez", Age: 28, Income: 42000, FamilySize: 1,
		Location: Urban, Segment: BudgetConscious, TechSavviness: TechIntermediate, PriceSensitivity: VeryHigh,
		Preferences: map[string]interface{}{"store_brand": true, "coupons": true},
		Behaviors:   []Behavior{{Type: WeeklyGrocery, Frequency: Weekly}},
	}
	df.personas["patricia_premium"] = &Persona{
		ID: "patricia_premium", Name: "Patricia Chen", Age: 45, Income: 285000, FamilySize: 4,
		Location: Suburban, Segment: PremiumSeeker, TechSavviness: TechAdvanced, PriceSensitivity: Low,
		Preferences: map[string]interface{}{"organic": true, "amazon_prime": true},
		Behaviors:   []Behavior{{Type: SubscriptionUser, Frequency: Monthly}},
	}
	df.personas["edward_elderly"] = &Persona{
		ID: "edward_elderly", Name: "Edward Thompson", Age: 72, Income: 38000, FamilySize: 2,
		Location: Rural, Segment: Elderly, TechSavviness: TechNovice, PriceSensitivity: Medium,
		Preferences: map[string]interface{}{"pharmacy_proximity": true},
		Behaviors:   []Behavior{{Type: WeeklyGrocery, Frequency: Weekly}},
	}
	df.personas["tyler_techsavvy"] = &Persona{
		ID: "tyler_techsavvy", Name: "Tyler Rodriguez", Age: 24, Income: 95000, FamilySize: 1,
		Location: Urban, Segment: TechSavvy, TechSavviness: TechExpert, PriceSensitivity: Medium,
		Preferences: map[string]interface{}{"app_only": true, "latest_tech": true},
		Behaviors:   []Behavior{{Type: ResearchIntensive, Frequency: Monthly}},
	}
	df.personas["fatima_family"] = &Persona{
		ID: "fatima_family", Name: "Fatima Johnson", Age: 36, Income: 78000, FamilySize: 5,
		Location: Suburban, Segment: Family, TechSavviness: TechIntermediate, PriceSensitivity: High,
		Preferences: map[string]interface{}{"bulk_buying": true},
		Behaviors:   []Behavior{{Type: WeeklyGrocery, Frequency: Weekly}},
	}
	df.personas["helen_health"] = &Persona{
		ID: "helen_health", Name: "Helen Park", Age: 52, Income: 125000, FamilySize: 2,
		Location: Suburban, Segment: HealthFocused, TechSavviness: TechAdvanced, PriceSensitivity: Low,
		Preferences: map[string]interface{}{"organic": true},
		Behaviors:   []Behavior{{Type: SubscriptionUser, Frequency: Monthly}},
	}
	df.personas["carlos_convenience"] = &Persona{
		ID: "carlos_convenience", Name: "Carlos Ramirez", Age: 41, Income: 165000, FamilySize: 3,
		Location: Urban, Segment: ConvenienceFirst, TechSavviness: TechAdvanced, PriceSensitivity: VeryLow,
		Preferences: map[string]interface{}{"same_day_delivery": true},
		Behaviors:   []Behavior{{Type: SubscriptionUser, Frequency: Monthly}},
	}
}

func (df *DataFactory) initializeProducts() {
	df.products[Walmart] = make(map[string]*Product)
	df.products[Amazon] = make(map[string]*Product)
	df.products[Apple] = make(map[string]*Product)
	df.products[CVS] = make(map[string]*Product)
	df.products[UnitedHealth] = make(map[string]*Product)
	df.initializeWalmartProducts()
	df.initializeAmazonProducts()
	df.initializeAppleProducts()
	df.initializeCVSProducts()
	df.initializeUnitedHealthProducts()
}

func (df *DataFactory) initializeWalmartProducts() {
	products := []Product{
		{ID: "WM-PROD-001", Name: "Organic Bananas", Category: "Produce", Brand: "Great Value", Price: 0.58, Company: Walmart, Availability: true, Rating: 4.5},
		{ID: "WM-DAIRY-001", Name: "Whole Milk", Category: "Dairy", Brand: "Great Value", Price: 3.87, Company: Walmart, Availability: true, Rating: 4.5},
		{ID: "WM-DAIRY-002", Name: "Large Eggs", Category: "Dairy", Brand: "Great Value", Price: 2.78, Company: Walmart, Availability: true, Rating: 4.6},
		{ID: "WM-MEAT-001", Name: "Chicken Breast", Category: "Meat", Brand: "All Natural", Price: 3.47, Company: Walmart, Availability: true, Rating: 4.4},
		{ID: "WM-PANTRY-001", Name: "White Bread", Category: "Bakery", Brand: "Great Value", Price: 1.18, Company: Walmart, Availability: true, Rating: 4.2},
		{ID: "WM-BEV-001", Name: "Cola 12-Pack", Category: "Beverages", Brand: "Coca-Cola", Price: 5.98, Company: Walmart, Availability: true, Rating: 4.7},
		{ID: "WM-HOUSE-001", Name: "Paper Towels", Category: "Household", Brand: "Bounty", Price: 24.98, Company: Walmart, Availability: true, Rating: 4.8},
	}
	for _, p := range products {
		product := p
		df.products[Walmart][product.ID] = &product
	}
}

func (df *DataFactory) initializeAmazonProducts() {
	products := []Product{
		{ID: "AMZ-ELEC-001", Name: "Echo Dot", Category: "Smart Home", Brand: "Amazon", Price: 49.99, Company: Amazon, Availability: true, Rating: 4.7},
		{ID: "AMZ-ELEC-002", Name: "Fire TV Stick", Category: "Streaming", Brand: "Amazon", Price: 49.99, Company: Amazon, Availability: true, Rating: 4.6},
		{ID: "AMZ-HOME-001", Name: "Instant Pot", Category: "Kitchen", Brand: "Instant Pot", Price: 89.00, Company: Amazon, Availability: true, Rating: 4.7},
		{ID: "AMZ-HEALTH-001", Name: "Vitamin D3", Category: "Vitamins", Brand: "Nature Made", Price: 14.99, Company: Amazon, Availability: true, Rating: 4.7},
	}
	for _, p := range products {
		product := p
		df.products[Amazon][product.ID] = &product
	}
}

func (df *DataFactory) initializeAppleProducts() {
	products := []Product{
		{ID: "AAPL-IP-001", Name: "iPhone 15 Pro Max", Category: "Smartphone", Brand: "Apple", Price: 1199.00, Company: Apple, Availability: true, Rating: 4.8},
		{ID: "AAPL-IPAD-001", Name: "iPad Pro", Category: "Tablet", Brand: "Apple", Price: 1199.00, Company: Apple, Availability: true, Rating: 4.9},
		{ID: "AAPL-MAC-001", Name: "MacBook Pro 16", Category: "Laptop", Brand: "Apple", Price: 3499.00, Company: Apple, Availability: true, Rating: 4.9},
		{ID: "AAPL-WATCH-001", Name: "Apple Watch Series 9", Category: "Smartwatch", Brand: "Apple", Price: 429.00, Company: Apple, Availability: true, Rating: 4.8},
		{ID: "AAPL-AIR-001", Name: "AirPods Pro", Category: "Audio", Brand: "Apple", Price: 249.00, Company: Apple, Availability: true, Rating: 4.8},
	}
	for _, p := range products {
		product := p
		df.products[Apple][product.ID] = &product
	}
}

func (df *DataFactory) initializeCVSProducts() {
	products := []Product{
		{ID: "CVS-RX-001", Name: "Atorvastatin 20mg", Category: "Prescription", Brand: "Generic", Price: 12.99, Company: CVS, Availability: true, Rating: 4.7},
		{ID: "CVS-RX-002", Name: "Metformin 500mg", Category: "Prescription", Brand: "Generic", Price: 8.99, Company: CVS, Availability: true, Rating: 4.6},
		{ID: "CVS-RX-003", Name: "Lisinopril 10mg", Category: "Prescription", Brand: "Generic", Price: 9.99, Company: CVS, Availability: true, Rating: 4.7},
		{ID: "CVS-OTC-001", Name: "Ibuprofen", Category: "Pain Relief", Brand: "CVS Health", Price: 8.99, Company: CVS, Availability: true, Rating: 4.5},
		{ID: "CVS-VIT-001", Name: "Multivitamin", Category: "Vitamins", Brand: "CVS Health", Price: 14.99, Company: CVS, Availability: true, Rating: 4.6},
		{ID: "CVS-VIT-002", Name: "Vitamin D3", Category: "Vitamins", Brand: "CVS Health", Price: 12.99, Company: CVS, Availability: true, Rating: 4.7},
		{ID: "CVS-MED-001", Name: "Bandages", Category: "First Aid", Brand: "CVS Health", Price: 4.99, Company: CVS, Availability: true, Rating: 4.5},
		{ID: "CVS-MED-002", Name: "Thermometer", Category: "Medical Devices", Brand: "CVS Health", Price: 9.99, Company: CVS, Availability: true, Rating: 4.6},
	}
	for _, p := range products {
		product := p
		df.products[CVS][product.ID] = &product
	}
}

func (df *DataFactory) initializeUnitedHealthProducts() {
	products := []Product{
		{ID: "UH-IND-001", Name: "Bronze Plan", Category: "Individual Health", Brand: "UnitedHealthcare", Price: 385.00, Company: UnitedHealth, Availability: true, Rating: 4.2},
		{ID: "UH-IND-002", Name: "Silver Plan", Category: "Individual Health", Brand: "UnitedHealthcare", Price: 525.00, Company: UnitedHealth, Availability: true, Rating: 4.4},
		{ID: "UH-IND-003", Name: "Gold Plan", Category: "Individual Health", Brand: "UnitedHealthcare", Price: 695.00, Company: UnitedHealth, Availability: true, Rating: 4.6},
		{ID: "UH-FAM-001", Name: "Family Silver", Category: "Family Health", Brand: "UnitedHealthcare", Price: 1450.00, Company: UnitedHealth, Availability: true, Rating: 4.5},
		{ID: "UH-FAM-002", Name: "Family Gold", Category: "Family Health", Brand: "UnitedHealthcare", Price: 1850.00, Company: UnitedHealth, Availability: true, Rating: 4.6},
		{ID: "UH-MED-001", Name: "Medicare Advantage HMO", Category: "Medicare", Brand: "UnitedHealthcare", Price: 0.00, Company: UnitedHealth, Availability: true, Rating: 4.5},
		{ID: "UH-DENT-001", Name: "Dental PPO", Category: "Dental", Brand: "UnitedHealthcare", Price: 35.00, Company: UnitedHealth, Availability: true, Rating: 4.4},
		{ID: "UH-VIS-001", Name: "Vision Plan", Category: "Vision", Brand: "UnitedHealthcare", Price: 12.00, Company: UnitedHealth, Availability: true, Rating: 4.3},
	}
	for _, p := range products {
		product := p
		df.products[UnitedHealth][product.ID] = &product
	}
}

type PersonaBuilder struct {
	persona *Persona
}

func NewPersonaBuilder(id, name string) *PersonaBuilder {
	return &PersonaBuilder{persona: &Persona{ID: id, Name: name, Preferences: make(map[string]interface{}), Behaviors: []Behavior{}}}
}

func (pb *PersonaBuilder) WithAge(age int) *PersonaBuilder                              { pb.persona.Age = age; return pb }
func (pb *PersonaBuilder) WithIncome(income int) *PersonaBuilder                        { pb.persona.Income = income; return pb }
func (pb *PersonaBuilder) WithFamilySize(size int) *PersonaBuilder                      { pb.persona.FamilySize = size; return pb }
func (pb *PersonaBuilder) WithLocation(loc LocationType) *PersonaBuilder                { pb.persona.Location = loc; return pb }
func (pb *PersonaBuilder) WithSegment(seg CustomerSegment) *PersonaBuilder              { pb.persona.Segment = seg; return pb }
func (pb *PersonaBuilder) WithTechSavviness(level TechLevel) *PersonaBuilder            { pb.persona.TechSavviness = level; return pb }
func (pb *PersonaBuilder) WithPriceSensitivity(sens PriceSensitivity) *PersonaBuilder   { pb.persona.PriceSensitivity = sens; return pb }
func (pb *PersonaBuilder) WithPreference(key string, value interface{}) *PersonaBuilder { pb.persona.Preferences[key] = value; return pb }
func (pb *PersonaBuilder) WithBehavior(b Behavior) *PersonaBuilder                      { pb.persona.Behaviors = append(pb.persona.Behaviors, b); return pb }
func (pb *PersonaBuilder) Build() *Persona                                              { return pb.persona }

type ScenarioBuilder struct {
	persona     *Persona
	context     *Context
	products    []*Product
	constraints Constraints
}

func NewScenarioBuilder() *ScenarioBuilder {
	return &ScenarioBuilder{context: &Context{}, products: []*Product{}, constraints: Constraints{Availability: make(map[string]bool)}}
}

func (sb *ScenarioBuilder) WithPersona(p *Persona) *ScenarioBuilder                 { sb.persona = p; return sb }
func (sb *ScenarioBuilder) WithTimeContext(tc TimeContext) *ScenarioBuilder         { sb.context.TimeContext = tc; return sb }
func (sb *ScenarioBuilder) WithLocationContext(lc LocationContext) *ScenarioBuilder { sb.context.LocationContext = lc; return sb }
func (sb *ScenarioBuilder) WithWeatherContext(wc WeatherContext) *ScenarioBuilder   { sb.context.WeatherContext = wc; return sb }
func (sb *ScenarioBuilder) WithEventContext(ec EventContext) *ScenarioBuilder       { sb.context.EventContext = ec; return sb }
func (sb *ScenarioBuilder) WithProducts(products ...*Product) *ScenarioBuilder      { sb.products = append(sb.products, products...); return sb }
func (sb *ScenarioBuilder) WithBudget(budget float64) *ScenarioBuilder              { sb.constraints.Budget = budget; return sb }
func (sb *ScenarioBuilder) WithTimeLimit(limit time.Duration) *ScenarioBuilder      { sb.constraints.TimeLimit = limit; return sb }
func (sb *ScenarioBuilder) Build() map[string]interface{} {
	sb.context.Constraints = sb.constraints
	return map[string]interface{}{"persona": sb.persona, "context": sb.context, "products": sb.products}
}

func (df *DataFactory) GetWalmartGroceryScenario(personaID string) map[string]interface{} {
	persona := df.personas[personaID]
	if persona == nil {
		persona = df.personas["sarah_budget"]
	}
	return NewScenarioBuilder().WithPersona(persona).WithTimeContext(Weekend).
		WithLocationContext(LocationContext{Type: Suburban, Distance: 2.5}).WithBudget(100.00).
		WithProducts(df.products[Walmart]["WM-PROD-001"], df.products[Walmart]["WM-DAIRY-001"]).Build()
}

func (df *DataFactory) GetAmazonPrimeScenario(personaID string) map[string]interface{} {
	persona := df.personas[personaID]
	if persona == nil {
		persona = df.personas["tyler_techsavvy"]
	}
	return NewScenarioBuilder().WithPersona(persona).WithTimeContext(LateNight).WithBudget(500.00).
		WithProducts(df.products[Amazon]["AMZ-ELEC-001"]).Build()
}

func (df *DataFactory) GetAppleEcosystemScenario(personaID string) map[string]interface{} {
	persona := df.personas[personaID]
	if persona == nil {
		persona = df.personas["patricia_premium"]
	}
	return NewScenarioBuilder().WithPersona(persona).WithEventContext(EventContext{Type: "product_launch", Urgency: "high"}).
		WithBudget(2000.00).WithProducts(df.products[Apple]["AAPL-IP-001"]).Build()
}

func (df *DataFactory) GetCVSPharmacyScenario(personaID string) map[string]interface{} {
	persona := df.personas[personaID]
	if persona == nil {
		persona = df.personas["edward_elderly"]
	}
	return NewScenarioBuilder().WithPersona(persona).WithEventContext(EventContext{Type: "prescription_refill"}).
		WithBudget(150.00).WithProducts(df.products[CVS]["CVS-RX-001"]).Build()
}

func (df *DataFactory) GetUnitedHealthEnrollmentScenario(personaID string) map[string]interface{} {
	persona := df.personas[personaID]
	if persona == nil {
		persona = df.personas["fatima_family"]
	}
	return NewScenarioBuilder().WithPersona(persona).WithTimeContext(HolidaySeason).
		WithEventContext(EventContext{Type: "open_enrollment", Urgency: "high"}).WithBudget(2500.00).
		WithProducts(df.products[UnitedHealth]["UH-FAM-002"]).Build()
}

func (df *DataFactory) GetPersona(id string) *Persona {
	return df.personas[id]
}

func (df *DataFactory) GetAllPersonas() map[string]*Persona {
	return df.personas
}

func (df *DataFactory) GetProduct(company Fortune5Company, productID string) *Product {
	if products, ok := df.products[company]; ok {
		return products[productID]
	}
	return nil
}

func (df *DataFactory) GetProductsByCompany(company Fortune5Company) map[string]*Product {
	return df.products[company]
}

func (df *DataFactory) GetProductsByCategory(company Fortune5Company, category string) []*Product {
	var result []*Product
	if products, ok := df.products[company]; ok {
		for _, product := range products {
			if product.Category == category {
				result = append(result, product)
			}
		}
	}
	return result
}

func (df *DataFactory) GenerateRandomTransaction(personaID string, company Fortune5Company, itemCount int) *Transaction {
	persona := df.personas[personaID]
	if persona == nil || itemCount <= 0 {
		return nil
	}
	companyProducts := df.products[company]
	if len(companyProducts) == 0 {
		return nil
	}

	var purchases []ProductPurchase
	total := 0.0
	productList := make([]*Product, 0, len(companyProducts))
	for _, p := range companyProducts {
		productList = append(productList, p)
	}

	for i := 0; i < itemCount && i < len(productList); i++ {
		idx := df.rand.Intn(len(productList))
		product := productList[idx]
		quantity := df.rand.Intn(3) + 1
		purchase := ProductPurchase{Product: product, Quantity: quantity, Price: product.Price * float64(quantity)}
		purchases = append(purchases, purchase)
		total += purchase.Price
	}

	return &Transaction{
		ID: fmt.Sprintf("TXN-%s-%d", personaID, time.Now().Unix()),
		PersonaID: personaID, Products: purchases, TotalAmount: total,
		Timestamp: time.Now(), Channel: Online, Context: &Context{},
	}
}

func (df *DataFactory) GenerateWeeklyGroceryList(personaID string) *Transaction {
	persona := df.personas[personaID]
	if persona == nil {
		return nil
	}
	groceryIDs := []string{"WM-PROD-001", "WM-DAIRY-001", "WM-DAIRY-002", "WM-MEAT-001"}
	var purchases []ProductPurchase
	total := 0.0
	for _, id := range groceryIDs {
		if product := df.products[Walmart][id]; product != nil {
			quantity := 1
			if persona.FamilySize > 2 {
				quantity = 2
			}
			purchase := ProductPurchase{Product: product, Quantity: quantity, Price: product.Price * float64(quantity)}
			purchases = append(purchases, purchase)
			total += purchase.Price
		}
	}
	return &Transaction{
		ID: fmt.Sprintf("TXN-GROCERY-%s-%d", personaID, time.Now().Unix()),
		PersonaID: personaID, Products: purchases, TotalAmount: total,
		Timestamp: time.Now(), Channel: InStore,
		Context: &Context{TimeContext: Weekend, LocationContext: LocationContext{Type: Suburban}},
	}
}

func (df *DataFactory) GetTestScenarios() []map[string]interface{} {
	return []map[string]interface{}{
		df.GetWalmartGroceryScenario("sarah_budget"),
		df.GetAmazonPrimeScenario("tyler_techsavvy"),
		df.GetAppleEcosystemScenario("patricia_premium"),
		df.GetCVSPharmacyScenario("edward_elderly"),
		df.GetUnitedHealthEnrollmentScenario("fatima_family"),
	}
}

func (df *DataFactory) PrintPersonaSummary(personaID string) string {
	persona := df.personas[personaID]
	if persona == nil {
		return "Persona not found"
	}
	return fmt.Sprintf("Persona: %s | Age: %d | Income: $%d | Segment: %s",
		persona.Name, persona.Age, persona.Income, persona.Segment)
}

func (df *DataFactory) GetStatistics() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["total_personas"] = len(df.personas)
	productCounts := make(map[Fortune5Company]int)
	for company, products := range df.products {
		productCounts[company] = len(products)
	}
	stats["products_by_company"] = productCounts
	return stats
}

func (df *DataFactory) ClonePersona(personaID string) *Persona {
	original := df.personas[personaID]
	if original == nil {
		return nil
	}
	clone := &Persona{
		ID: original.ID + "_clone", Name: original.Name, Age: original.Age, Income: original.Income,
		FamilySize: original.FamilySize, Location: original.Location, Segment: original.Segment,
		TechSavviness: original.TechSavviness, PriceSensitivity: original.PriceSensitivity,
		Preferences: make(map[string]interface{}), Behaviors: make([]Behavior, len(original.Behaviors)),
	}
	for k, v := range original.Preferences {
		clone.Preferences[k] = v
	}
	copy(clone.Behaviors, original.Behaviors)
	return clone
}
