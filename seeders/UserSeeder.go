package seeders

import (
	"log"

	"gorm.io/gorm"
)

// UserSeeder is deprecated: user identities and memberships are managed centrally by Central SSO.
func UserSeeder(_ *gorm.DB) error {
	log.Println("✓ UserSeeder skipped: user accounts are owned exclusively by Central SSO")
	return nil
}
