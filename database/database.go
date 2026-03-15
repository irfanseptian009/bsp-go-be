package database

import (
	"fmt"
	"log"

	"github.com/irfanseptian/fims-backend/config"
	"github.com/irfanseptian/fims-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance.
var DB *gorm.DB

// Connect initializes the database connection and runs auto-migration.
func Connect(cfg *config.Config) {
	var err error

	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connected successfully")

	// Ensure UUID extension and DB-level defaults are present.
	ensureUUIDDefaults()

	// Drop legacy FK constraints (if any) to prevent type mismatch migration failures.
	dropLegacyConstraints()

	// Auto-migrate all models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Branch{},
		&models.OccupationType{},
		&models.InsuranceRequest{},
		&models.Policy{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to auto-migrate: %v", err)
	}

	log.Println("✅ Database migration completed")

	// Create enum types if not exists (PostgreSQL)
	createEnumTypes()
}

// createEnumTypes ensures PostgreSQL enum types exist.
func createEnumTypes() {
	enums := []struct {
		typeName string
		values   string
	}{
		{"role", "'CUSTOMER', 'ADMIN'"},
		{"request_status", "'PENDING', 'APPROVED', 'REJECTED'"},
		{"construction_class", "'KELAS_1', 'KELAS_2', 'KELAS_3'"},
	}

	for _, e := range enums {
		sql := fmt.Sprintf(
			"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN CREATE TYPE %s AS ENUM (%s); END IF; END $$;",
			e.typeName, e.typeName, e.values,
		)
		DB.Exec(sql)
	}
}

func dropLegacyConstraints() {
	queries := []string{
		`ALTER TABLE IF EXISTS insurance_requests DROP CONSTRAINT IF EXISTS insurance_requests_user_id_fkey;`,
		`ALTER TABLE IF EXISTS insurance_requests DROP CONSTRAINT IF EXISTS fk_insurance_requests_user;`,
		`ALTER TABLE IF EXISTS insurance_requests DROP CONSTRAINT IF EXISTS insurance_requests_occupation_type_id_fkey;`,
		`ALTER TABLE IF EXISTS insurance_requests DROP CONSTRAINT IF EXISTS fk_insurance_requests_occupation_type;`,
		`ALTER TABLE IF EXISTS policies DROP CONSTRAINT IF EXISTS policies_branch_id_fkey;`,
		`ALTER TABLE IF EXISTS policies DROP CONSTRAINT IF EXISTS fk_policies_branch;`,
		`ALTER TABLE IF EXISTS policies DROP CONSTRAINT IF EXISTS policies_occupation_type_id_fkey;`,
		`ALTER TABLE IF EXISTS policies DROP CONSTRAINT IF EXISTS fk_policies_occupation_type;`,
	}

	for _, q := range queries {
		if err := DB.Exec(q).Error; err != nil {
			log.Printf("⚠️ Failed to drop legacy constraint: %v", err)
		}
	}
}

func ensureUUIDDefaults() {
	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS pgcrypto;`,
		`ALTER TABLE IF EXISTS users ALTER COLUMN id SET DEFAULT gen_random_uuid();`,
		`ALTER TABLE IF EXISTS branches ALTER COLUMN id SET DEFAULT gen_random_uuid();`,
		`ALTER TABLE IF EXISTS occupation_types ALTER COLUMN id SET DEFAULT gen_random_uuid();`,
		`ALTER TABLE IF EXISTS insurance_requests ALTER COLUMN id SET DEFAULT gen_random_uuid();`,
		`ALTER TABLE IF EXISTS policies ALTER COLUMN id SET DEFAULT gen_random_uuid();`,
	}

	for _, q := range queries {
		if err := DB.Exec(q).Error; err != nil {
			log.Printf("⚠️ Failed to ensure UUID default: %v", err)
		}
	}
}

// GetDB returns the database instance.
func GetDB() *gorm.DB {
	return DB
}
