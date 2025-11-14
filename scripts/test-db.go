package main

import (
	"fmt"
	"log"
)

type Country struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Code string
}

type Technology struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	CountryID    uint
	Manufacturer string
	Status       string
}

func main() {
	fmt.Println("🔍 Testing database connection...")

	// Connect to database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("❌ Connection failed:", err)
	}

	fmt.Println("✅ Database connection successful!")

	// Health check
	if err := config.HealthCheck(db); err != nil {
		log.Fatal("❌ Health check failed:", err)
	}
	fmt.Println("✅ Health check passed!")

	// Test query - countries
	var countries []Country
	result := db.Find(&countries)
	if result.Error != nil {
		log.Fatal("❌ Query failed:", result.Error)
	}

	fmt.Printf("\n📊 Found %d countries:\n", len(countries))
	for _, c := range countries {
		fmt.Printf("  • %s (%s)\n", c.Name, c.Code)
	}

	// Test query - technologies
	var techs []Technology
	result = db.Limit(5).Find(&techs)
	if result.Error != nil {
		log.Fatal("❌ Query failed:", result.Error)
	}

	fmt.Printf("\n🚀 Sample technologies (%d total):\n", len(techs))
	for _, t := range techs {
		fmt.Printf("  • %s - %s (%s)\n", t.Name, t.Manufacturer, t.Status)
	}

	fmt.Println("\n✅ All tests passed!")
}
