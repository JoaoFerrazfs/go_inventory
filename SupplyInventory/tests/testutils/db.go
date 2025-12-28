package testutils

import (
	"log"
	"os"

	dbInfra "go_inventory/SupplyInventory/Infrastructure/Db"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func SetupTestDB() *gorm.DB {
	if TestDB != nil {
		return TestDB
	}
	var dsn string
	var dbName string = "inventory_test"

	dbHost := os.Getenv("TEST_DB_HOST")
	dbPort := os.Getenv("TEST_DB_PORT")
	dbUser := os.Getenv("TEST_DB_USER")
	dbPass := os.Getenv("TEST_DB_PASSWORD")

	if dbHost != "" && dbPort != "" && dbUser != "" {
		// First, connect without DB to create it
		dsnNoDB := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/?charset=utf8mb4&parseTime=True&loc=Local"
		dbTemp, err := gorm.Open(mysql.Open(dsnNoDB), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to DB server:", err)
		}
		sqlDBTemp, _ := dbTemp.DB()
		defer sqlDBTemp.Close()
		if err := sqlDBTemp.Ping(); err != nil {
			log.Fatal("Error pinging DB server:", err)
		}
		// Create DB if not exists
		dbTemp.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
		dsn = dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	} else {
		// fallback for local tests
		dsnNoDB := "root:root@tcp(db:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
		dbTemp, err := gorm.Open(mysql.Open(dsnNoDB), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to DB server:", err)
		}
		sqlDBTemp, _ := dbTemp.DB()
		defer sqlDBTemp.Close()
		if err := sqlDBTemp.Ping(); err != nil {
			log.Fatal("Error pinging DB server:", err)
		}
		// Create DB if not exists
		dbTemp.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
		dsn = "root:root@tcp(db:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test DB:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Error pinging test DB:", err)
	}

	log.Println("Test DB connection successful")

	// Run AutoMigrate
	if err := dbInfra.Migrate(db); err != nil {
		log.Fatal("Failed to migrate test DB:", err)
	}

	TestDB = db
	return db
}

func TeardownTestDB() {
	if TestDB != nil {
		sqlDB, _ := TestDB.DB()
		sqlDB.Close()
	}
}

func TruncateTables(db *gorm.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE TABLE user_entities")
	db.Exec("TRUNCATE TABLE pallet_racks")
	db.Exec("TRUNCATE TABLE pallets")
	db.Exec("TRUNCATE TABLE palletized_products")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}
