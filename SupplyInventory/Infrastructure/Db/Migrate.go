package db

import (
	"log"

	entities "go_inventory/SupplyInventory/Domain/Entities"

	"gorm.io/gorm"
)

// Migrate runs all database migrations
// Note: With RBAC implementation, UserEntity now includes Role field
// If you have existing database, manually add column:
// ALTER TABLE user_entities ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' AFTER password;
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entities.InventoryEntity{},
		&entities.PalletRackEntity{},
		&entities.PalletEntity{},
		&entities.PalletizedProductEntity{},
		&entities.UserEntity{},
		&entities.ProductEntity{},
	)
	if err != nil {
		return err
	}

	LogMigrationInstructions()
	return nil
}

// LogMigrationInstructions prints helpful information about RBAC migration
func LogMigrationInstructions() {
	log.Println(`
╔════════════════════════════════════════════════════════════════╗
║           RBAC (Role-Based Access Control) ENABLED            ║
╠════════════════════════════════════════════════════════════════╣
║                                                                ║
║ UserEntity now includes role-based access control!           ║
║                                                                ║
║ Roles Available:                                              ║
║   - "admin" : Full system access                             ║
║   - "user"  : Regular user (default)                         ║
║                                                                ║
║ If you have an EXISTING database:                            ║
║                                                                ║
║ 1. Add role column:                                           ║
║    ALTER TABLE user_entities                                  ║
║    ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user'       ║
║    AFTER password;                                            ║
║                                                                ║
║ 2. (Optional) Make first user admin:                          ║
║    UPDATE user_entities SET role = 'admin' WHERE id = 1;     ║
║                                                                ║
║ 3. Verify:                                                    ║
║    SELECT id, email, role FROM user_entities;                ║
║                                                                ║
║ See: documents/security/RBAC_GUIDE.md for complete guide    ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝
	`)
}
