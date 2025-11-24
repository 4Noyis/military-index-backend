package main

import (
	"context"
	"fmt"
	"log"

	"github.com/4Noyis/military-index-backend/services/technology-service/internal/repository"
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
	repo := repository.NewTechnologyRepository(db)
	ctx := context.Background()

	// Test 1: GetAll with pagination
	fmt.Println("--- Test 1: GetAll (page=1, limit=5) ---")
	technologies, total, err := repo.GetAll(ctx, 1, 5)
	if err != nil {
		log.Printf("❌ GetAll failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d technologies (total: %d)\n", len(technologies), total)
		for _, c := range technologies {
			fmt.Printf("   - %s | %s | %s\n", c.Name, c.Description, c.Category.Name)
		}
	}
	fmt.Println()

	fmt.Println("--- Test 2: GetByID (ID=1) ---")
	technology, err := repo.GetByID(ctx, 1)
	if err != nil {
		log.Printf("❌ GetByID failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found: %s | %s | %s\n", technology.Name, technology.Description, technology.ImageURL)
	}
	fmt.Println()

	fmt.Println("--- Test 3: GetByCountry (ID=TUR) ---")
	technologiesByCountry, total, err := repo.GetByCountry(ctx, "TUR")
	if err != nil {
		log.Printf("❌ GetAll failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d technologies (total: %d)\n", len(technologiesByCountry), total)
		for _, c := range technologiesByCountry {
			fmt.Printf("   - %s (%s) %s\n", c.Name, c.Description, c.Country.Code)
		}
	}
	fmt.Println()

	fmt.Println("--- Test 4: GetByCountry (Category: Aircraft) ---")
	technologiesByCategory, total, err := repo.GetByCategory(ctx, "Aircraft")
	if err != nil {
		log.Printf("❌ GetByCountry failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d technologies (total: %d)\n", len(technologiesByCategory), total)
		for _, c := range technologiesByCategory {
			fmt.Printf("   - %s | %s | %s | %s\n", c.Name, c.Description, c.Country.Code, c.Category.Name)
		}
	}
	fmt.Println()

	fmt.Println("--- Test 5: GetByStatus (status = current) ---")
	technologiesByStatus, total, err := repo.GetByStatus(ctx, "current")
	if err != nil {
		log.Printf("❌ GetByStatus failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d technologies (total: %d)\n", len(technologiesByStatus), total)
		for _, c := range technologiesByStatus {
			fmt.Printf("   - %s | %s | %s\n", c.Name, c.Description, c.Status)
		}
	}
	fmt.Println()

	fmt.Println("--- Test 6: GetByYearRange (2002-2025) ---")
	technologiesByYearRange, total, err := repo.GetByYearRange(ctx, 2002, 2025)
	if err != nil {
		log.Printf("❌ GetByYearRange failed: %v\n", err)
	} else {
		fmt.Printf("✅ Found %d technologies (total: %d)\n", len(technologiesByYearRange), total)
		for _, c := range technologiesByYearRange {
			if c.YearDeveloped != nil {
				fmt.Printf("   - %s | %s | %d\n", c.Name, c.Description, *c.YearDeveloped)
			} else {
				fmt.Printf("   - %s | %s | N/A\n", c.Name, c.Description)
			}
		}
	}
	fmt.Println()

	// Test 7: Create
	fmt.Println("--- Test 7: Create ---")
	yearDeveloped := 2020
	newTech := &models.Technology{
		CountryID:     1, // Turkey
		CategoryID:    1, // Aircraft
		Name:          "Test Technology",
		Description:   "A test technology for repository testing",
		Designer:      "Test Designer",
		YearDeveloped: &yearDeveloped,
		Manufacturer:  "Test Manufacturer",
		Status:        "prototype",
	}
	err = repo.Create(ctx, newTech)
	if err != nil {
		log.Printf("❌ Create failed: %v\n", err)
	} else {
		fmt.Printf("✅ Created technology with ID: %d\n", newTech.ID)
	}
	fmt.Println()

	// Test 8: Update
	fmt.Println("--- Test 8: Update ---")
	if newTech.ID > 0 {
		newTech.Name = "Updated Test Technology"
		newTech.Description = "Updated description"
		yearDeployed := 2023
		newTech.YearDeployed = &yearDeployed
		newTech.Status = "current"

		err = repo.Update(ctx, newTech)
		if err != nil {
			log.Printf("❌ Update failed: %v\n", err)
		} else {
			fmt.Printf("✅ Updated technology ID: %d\n", newTech.ID)
			// Verify update by fetching
			updated, err := repo.GetByID(ctx, newTech.ID)
			if err != nil {
				log.Printf("❌ Failed to verify update: %v\n", err)
			} else {
				fmt.Printf("   - Name: %s\n", updated.Name)
				fmt.Printf("   - Description: %s\n", updated.Description)
				fmt.Printf("   - Status: %s\n", updated.Status)
				if updated.YearDeployed != nil {
					fmt.Printf("   - Year Deployed: %d\n", *updated.YearDeployed)
				}
			}
		}
	} else {
		fmt.Println("⚠️  Skipping Update test - no technology was created")
	}
	fmt.Println()

	// Test 9: Delete
	fmt.Println("--- Test 9: Delete ---")
	if newTech.ID > 0 {
		err = repo.Delete(ctx, newTech.ID)
		if err != nil {
			log.Printf("❌ Delete failed: %v\n", err)
		} else {
			fmt.Printf("✅ Deleted technology ID: %d\n", newTech.ID)
			// Verify deletion
			_, err := repo.GetByID(ctx, newTech.ID)
			if err != nil {
				fmt.Printf("   - Verified: technology no longer exists\n")
			} else {
				fmt.Printf("   - ⚠️  Technology still exists after deletion\n")
			}
		}
	} else {
		fmt.Println("⚠️  Skipping Delete test - no technology was created")
	}
	fmt.Println()

}
