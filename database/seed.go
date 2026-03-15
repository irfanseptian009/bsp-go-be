package database

import (
	"log"

	"github.com/irfanseptian/fims-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed populates the database with initial data.
func Seed() {
	db := GetDB()

	log.Println("🌱 Seeding database...")

	seedBranches(db)
	seedOccupationTypes(db)
	seedUsers(db)

	log.Println("🎉 Seeding complete!")
	log.Println("─────────────────────────────────────")
	log.Println("📋 Test accounts:")
	log.Println("   Admin    : admin@fims.com     / admin123")
	log.Println("   Customer : customer@fims.com  / customer123")
	log.Println("   Customer : siti@fims.com      / customer123")
	log.Println("   Customer : andi@fims.com      / customer123")
	log.Println("─────────────────────────────────────")
}

func seedBranches(db *gorm.DB) {
	branches := []models.Branch{
		{Code: "00001", Name: "Kuningan"},
		{Code: "00002", Name: "Tebet"},
		{Code: "00003", Name: "Harmoni"},
		{Code: "00004", Name: "Sudirman"},
		{Code: "00005", Name: "Kelapa Gading"},
	}

	for _, b := range branches {
		var existing models.Branch
		result := db.Where("code = ?", b.Code).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			db.Create(&b)
		}
	}
	log.Println("✅ Branches seeded")
}

func seedOccupationTypes(db *gorm.DB) {
	types := []models.OccupationType{
		{Code: "2976.01", Name: "Rumah", PremiumRate: 0.3875},
		{Code: "2974.00", Name: "Ruko", PremiumRate: 0.5},
		{Code: "2975.00", Name: "Gedung Kantor", PremiumRate: 0.6},
		{Code: "2973.00", Name: "Gudang", PremiumRate: 0.75},
		{Code: "2972.00", Name: "Apartemen", PremiumRate: 0.45},
	}

	for _, t := range types {
		var existing models.OccupationType
		result := db.Where("code = ?", t.Code).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			db.Create(&t)
		} else {
			// Update premium rate if changed
			db.Model(&existing).Update("premium_rate", t.PremiumRate)
		}
	}
	log.Println("✅ Occupation types seeded")
}

func seedUsers(db *gorm.DB) {
	// Admin
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
	var adminUser models.User
	result := db.Where("email = ?", "admin@fims.com").First(&adminUser)
	if result.Error == gorm.ErrRecordNotFound {
		db.Create(&models.User{
			Name:     "Administrator",
			Email:    "admin@fims.com",
			Password: string(adminPassword),
			Role:     models.RoleAdmin,
		})
	}

	// Customers
	customerPassword, _ := bcrypt.GenerateFromPassword([]byte("customer123"), 10)
	customers := []models.User{
		{Name: "Budi Santoso", Email: "customer@fims.com", Password: string(customerPassword), Role: models.RoleCustomer},
		{Name: "Siti Rahayu", Email: "siti@fims.com", Password: string(customerPassword), Role: models.RoleCustomer},
		{Name: "Andi Wijaya", Email: "andi@fims.com", Password: string(customerPassword), Role: models.RoleCustomer},
	}

	for _, c := range customers {
		var existing models.User
		result := db.Where("email = ?", c.Email).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			db.Create(&c)
		}
	}
	log.Println("✅ Users seeded")
}
