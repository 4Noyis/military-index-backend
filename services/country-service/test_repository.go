package main

import (
	"context"
	"fmt"
	"log"

	"github.com/4Noyis/military-index-backend/services/country-service/internal/repository"
	"github.com/4Noyis/military-index-backend/shared/config"
	"github.com/4Noyis/military-index-backend/shared/models"
)

func main() {
	fmt.Println("=== Testing Country Repository ===\n")

	// Connect to database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("✅ Database connected\n")

	// Create repository
	repo := repository.NewCountryRepository(db)
	ctx := context.Background()

	// Test 1: GetAll with pagination
	fmt.Println("--- Test 1: GetAll (page=1, limit=5) ---")
	countries, total, err := repo.GetAll(ctx, 1, 5)
	if err != nil {
		log.Printf("❌ GetAll failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d countries (total: %d)\n", len(countries), total)
		for _, c := range countries {
			fmt.Printf("   - %s (%s)\n", c.Name, c.Code)
		}
	}
	fmt.Println()

	// Test 2: GetByID
	fmt.Println("--- Test 2: GetByID (ID=1) ---")
	country, err := repo.GetByID(ctx, 1)
	if err != nil {
		log.Printf("❌ GetByID failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found: %s (%s) - ID: %d\n", country.Name, country.Code, country.ID)
	}
	fmt.Println()

	// Test 3: GetByCode
	fmt.Println("--- Test 3: GetByCode (code='TUR') ---")
	country, err = repo.GetByCode(ctx, "TUR")
	if err != nil {
		log.Printf("❌ GetByCode failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found: %s (%s) - ID: %d\n", country.Name, country.Code, country.ID)
	}
	fmt.Println()

	// Test 4: Search
	fmt.Println("--- Test 4: Search (query='united') ---")
	countries, total, err = repo.Search(ctx, "united", 1, 10)
	if err != nil {
		log.Printf("❌ Search failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d countries matching 'united'\n", len(countries))
		for _, c := range countries {
			fmt.Printf("   - %s (%s)\n", c.Name, c.Code)
		}
	}
	fmt.Println()

	// Test 5: Create
	fmt.Println("--- Test 5: Create new country ---")
	newCountry := &models.Country{
		Name:    "Test Country",
		Code:    "TST",
		FlagURL: "/flags/test.svg",
	}
	err = repo.Create(ctx, newCountry)
	if err != nil {
		log.Printf("❌ Create failed: %v\n", err)
	} else {
		fmt.Printf("✅ Created: %s (%s) - ID: %d\n", newCountry.Name, newCountry.Code, newCountry.ID)
	}
	fmt.Println()

	// Test 6: Update
	if newCountry.ID > 0 {
		fmt.Println("--- Test 6: Update country ---")
		newCountry.Name = "Updated Test Country"
		err = repo.Update(ctx, newCountry)
		if err != nil {
			log.Printf("❌ Update failed: %v\n", err)
		} else {
			fmt.Printf("✅ Updated: %s (%s) - ID: %d\n", newCountry.Name, newCountry.Code, newCountry.ID)
		}
		fmt.Println()

		// Test 7: Delete
		fmt.Println("--- Test 7: Delete country ---")
		err = repo.Delete(ctx, newCountry.ID)
		if err != nil {
			log.Printf("❌ Delete failed: %v\n", err)
		} else {
			fmt.Printf("✅ Deleted country ID: %d\n", newCountry.ID)
		}
		fmt.Println()
	}

	// Test 8: GetByID with non-existent ID
	fmt.Println("--- Test 8: GetByID with non-existent ID (999999) ---")
	_, err = repo.GetByID(ctx, 999999)
	if err != nil {
		fmt.Printf("✅ Correctly returned error: %v\n", err)
	} else {
		fmt.Println("❌ Should have returned error for non-existent ID")
	}
	fmt.Println()

	// Test 9: Delete non-existent country
	fmt.Println("--- Test 9: Delete non-existent country (999999) ---")
	err = repo.Delete(ctx, 999999)
	if err != nil {
		fmt.Printf("✅ Correctly returned error: %v\n", err)
	} else {
		fmt.Println("❌ Should have returned error for non-existent ID")
	}
	fmt.Println()

	fmt.Println("=== All Repository Tests Complete ===")
}
