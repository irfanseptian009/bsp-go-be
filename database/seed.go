package database

import (
	"log"
	"time"

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
	seedInsuranceRequests(db)
	seedPolicies(db)

	log.Println("🎉 Seeding complete!")
	log.Println("─────────────────────────────────────")
	log.Println("📋 Test accounts:")
	log.Println("   Admin    : admin@bsp.com     / admin123")
	log.Println("   Customer : customer@bsp.com  / customer123")
	log.Println("   Customer : siti@bsp.com      / customer123")
	log.Println("   Customer : andi@bsp.com      / customer123")
	log.Println("   Customer : dewi@bsp.com      / customer123")
	log.Println("   Customer : rizky@bsp.com     / customer123")
	log.Println("─────────────────────────────────────")
}

// Reset truncates application tables.
func Reset() {
	db := GetDB()

	log.Println("🧹 Resetting database...")

	if err := db.Exec(`
		TRUNCATE TABLE
			policies,
			insurance_requests,
			users,
			branches,
			occupation_types
		RESTART IDENTITY CASCADE;
	`).Error; err != nil {
		log.Fatalf("❌ Failed to reset database: %v", err)
	}

	log.Println("✅ Database reset complete")
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
	type userSeed struct {
		Name     string
		Email    string
		Password string
		Role     models.Role
	}

	users := []userSeed{
		{Name: "Administrator", Email: "admin@bsp.com", Password: "admin123", Role: models.RoleAdmin},
		{Name: "Budi Santoso", Email: "customer@bsp.com", Password: "customer123", Role: models.RoleCustomer},
		{Name: "Siti Rahayu", Email: "siti@bsp.com", Password: "customer123", Role: models.RoleCustomer},
		{Name: "Andi Wijaya", Email: "andi@bsp.com", Password: "customer123", Role: models.RoleCustomer},
		{Name: "Dewi Lestari", Email: "dewi@bsp.com", Password: "customer123", Role: models.RoleCustomer},
		{Name: "Rizky Pratama", Email: "rizky@bsp.com", Password: "customer123", Role: models.RoleCustomer},
	}

	for _, u := range users {
		var existing models.User
		result := db.Where("email = ?", u.Email).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
			if err != nil {
				log.Printf("⚠️ Failed to hash password for %s: %v", u.Email, err)
				continue
			}

			db.Create(&models.User{
				Name:     u.Name,
				Email:    u.Email,
				Password: string(hashedPassword),
				Role:     u.Role,
			})
		}
	}
	log.Println("✅ Users seeded")
}

func seedInsuranceRequests(db *gorm.DB) {
	type requestSeed struct {
		InvoiceNumber     string
		UserEmail         string
		OccupationCode    string
		BuildingPrice     float64
		Duration          int
		ConstructionClass models.ConstructionClass
		Address           string
		Province          string
		City              string
		District          string
		Area              string
		Earthquake        bool
		Status            models.RequestStatus
		PolicyNumber      *string
	}

	requests := []requestSeed{
		{
			InvoiceNumber:     "K.001.10001",
			UserEmail:         "customer@bsp.com",
			OccupationCode:    "2976.01",
			BuildingPrice:     750000000,
			Duration:          1,
			ConstructionClass: models.ConstructionKelas1,
			Address:           "Jl. Melati No. 10",
			Province:          "DKI Jakarta",
			City:              "Jakarta Selatan",
			District:          "Kebayoran Baru",
			Area:              "Gandaria",
			Earthquake:        false,
			Status:            models.StatusApproved,
			PolicyNumber:      stringPtr("K.01.001.20001"),
		},
		{
			InvoiceNumber:     "K.001.10002",
			UserEmail:         "siti@bsp.com",
			OccupationCode:    "2972.00",
			BuildingPrice:     650000000,
			Duration:          2,
			ConstructionClass: models.ConstructionKelas2,
			Address:           "Jl. Anggrek No. 22",
			Province:          "Jawa Barat",
			City:              "Bandung",
			District:          "Coblong",
			Area:              "Dago",
			Earthquake:        true,
			Status:            models.StatusPending,
		},
		{
			InvoiceNumber:     "K.001.10003",
			UserEmail:         "andi@bsp.com",
			OccupationCode:    "2974.00",
			BuildingPrice:     1200000000,
			Duration:          1,
			ConstructionClass: models.ConstructionKelas2,
			Address:           "Jl. Veteran No. 8",
			Province:          "Jawa Timur",
			City:              "Surabaya",
			District:          "Genteng",
			Area:              "Ketabang",
			Earthquake:        false,
			Status:            models.StatusRejected,
		},
		{
			InvoiceNumber:     "K.001.10004",
			UserEmail:         "dewi@bsp.com",
			OccupationCode:    "2975.00",
			BuildingPrice:     2300000000,
			Duration:          3,
			ConstructionClass: models.ConstructionKelas1,
			Address:           "Jl. Pemuda No. 15",
			Province:          "Jawa Tengah",
			City:              "Semarang",
			District:          "Semarang Tengah",
			Area:              "Pandanaran",
			Earthquake:        true,
			Status:            models.StatusApproved,
			PolicyNumber:      stringPtr("K.01.001.20002"),
		},
		{
			InvoiceNumber:     "K.001.10005",
			UserEmail:         "rizky@bsp.com",
			OccupationCode:    "2973.00",
			BuildingPrice:     1800000000,
			Duration:          2,
			ConstructionClass: models.ConstructionKelas3,
			Address:           "Jl. Sudirman No. 99",
			Province:          "Banten",
			City:              "Tangerang",
			District:          "Cipondoh",
			Area:              "Poris",
			Earthquake:        false,
			Status:            models.StatusPending,
		},
	}

	userEmails := []string{"customer@bsp.com", "siti@bsp.com", "andi@bsp.com", "dewi@bsp.com", "rizky@bsp.com"}
	var users []models.User
	db.Where("email IN ?", userEmails).Find(&users)
	userIDByEmail := map[string]string{}
	for _, u := range users {
		userIDByEmail[u.Email] = u.ID
	}

	occupationCodes := []string{"2976.01", "2972.00", "2974.00", "2975.00", "2973.00"}
	var occupationTypes []models.OccupationType
	db.Where("code IN ?", occupationCodes).Find(&occupationTypes)
	occupationByCode := map[string]models.OccupationType{}
	for _, ot := range occupationTypes {
		occupationByCode[ot.Code] = ot
	}

	for _, r := range requests {
		userID, ok := userIDByEmail[r.UserEmail]
		if !ok {
			log.Printf("⚠️ User not found for request seed: %s", r.UserEmail)
			continue
		}

		occupationType, ok := occupationByCode[r.OccupationCode]
		if !ok {
			log.Printf("⚠️ Occupation type not found for request seed: %s", r.OccupationCode)
			continue
		}

		basicPremium := (r.BuildingPrice * occupationType.PremiumRate) / 1000 * float64(r.Duration)
		totalAmount := basicPremium + 10000

		var existing models.InsuranceRequest
		result := db.Where("invoice_number = ?", r.InvoiceNumber).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			db.Create(&models.InsuranceRequest{
				UserID:            userID,
				InvoiceNumber:     r.InvoiceNumber,
				OccupationTypeID:  occupationType.ID,
				BuildingPrice:     r.BuildingPrice,
				Duration:          r.Duration,
				ConstructionClass: r.ConstructionClass,
				Address:           r.Address,
				Province:          r.Province,
				City:              r.City,
				District:          r.District,
				Area:              r.Area,
				Earthquake:        r.Earthquake,
				BasicPremium:      basicPremium,
				AdminFee:          10000,
				TotalAmount:       totalAmount,
				Status:            r.Status,
				PolicyNumber:      r.PolicyNumber,
			})
		}
	}

	log.Println("✅ Insurance requests seeded")
}

func seedPolicies(db *gorm.DB) {
	type policySeed struct {
		PolicyNumber      string
		ApplicationNumber string
		Name              string
		BranchCode        string
		BirthDate         time.Time
		Duration          int
		BuildingPrice     float64
		OccupationCode    string
		Premium           float64
	}

	policies := []policySeed{
		{
			PolicyNumber:      "PL-2026-0001",
			ApplicationNumber: "APP-2026-0001",
			Name:              "Budi Santoso",
			BranchCode:        "00001",
			BirthDate:         time.Date(1990, 3, 10, 0, 0, 0, 0, time.UTC),
			Duration:          1,
			BuildingPrice:     750000000,
			OccupationCode:    "2976.01",
			Premium:           290625,
		},
		{
			PolicyNumber:      "PL-2026-0002",
			ApplicationNumber: "APP-2026-0002",
			Name:              "Dewi Lestari",
			BranchCode:        "00004",
			BirthDate:         time.Date(1988, 7, 21, 0, 0, 0, 0, time.UTC),
			Duration:          3,
			BuildingPrice:     2300000000,
			OccupationCode:    "2975.00",
			Premium:           4140000,
		},
	}

	branchCodes := []string{"00001", "00004"}
	var branches []models.Branch
	db.Where("code IN ?", branchCodes).Find(&branches)
	branchIDByCode := map[string]string{}
	for _, b := range branches {
		branchIDByCode[b.Code] = b.ID
	}

	occupationCodes := []string{"2976.01", "2975.00"}
	var occupationTypes []models.OccupationType
	db.Where("code IN ?", occupationCodes).Find(&occupationTypes)
	occupationIDByCode := map[string]string{}
	for _, ot := range occupationTypes {
		occupationIDByCode[ot.Code] = ot.ID
	}

	for _, p := range policies {
		branchID, ok := branchIDByCode[p.BranchCode]
		if !ok {
			log.Printf("⚠️ Branch not found for policy seed: %s", p.BranchCode)
			continue
		}

		occupationID, ok := occupationIDByCode[p.OccupationCode]
		if !ok {
			log.Printf("⚠️ Occupation type not found for policy seed: %s", p.OccupationCode)
			continue
		}

		var existing models.Policy
		result := db.Where("policy_number = ?", p.PolicyNumber).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			db.Create(&models.Policy{
				PolicyNumber:      p.PolicyNumber,
				ApplicationNumber: p.ApplicationNumber,
				Name:              p.Name,
				BranchID:          branchID,
				BirthDate:         p.BirthDate,
				Duration:          p.Duration,
				BuildingPrice:     p.BuildingPrice,
				OccupationTypeID:  occupationID,
				Premium:           p.Premium,
			})
		}
	}

	log.Println("✅ Policies seeded")
}

func stringPtr(v string) *string {
	return &v
}
