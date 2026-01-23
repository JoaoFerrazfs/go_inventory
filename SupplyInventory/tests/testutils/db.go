package testutils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	dbInfra "go_inventory/SupplyInventory/Infrastructure/Db"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	TestDB     *gorm.DB
	TestDBName string
)

func SetupTestDB() *gorm.DB {
	if TestDB != nil {
		return TestDB
	}
	var dsn string
	var dbName string = os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = fmt.Sprintf("inventory_test_%d", time.Now().UnixNano())
	}
	TestDBName = dbName

	dbHost := os.Getenv("TEST_DB_HOST")
	dbPort := os.Getenv("TEST_DB_PORT")
	dbUser := os.Getenv("TEST_DB_USER")
	dbPass := os.Getenv("TEST_DB_PASSWORD")

	if dbHost != "" && dbPort != "" && dbUser != "" {
		// Wait for DB server to be available (useful in CI where DB may start slowly)
		waitSeconds := 30
		if v := os.Getenv("TEST_DB_WAIT_SECONDS"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				waitSeconds = n
			}
		}
		retryIntervalMs := 500
		if v := os.Getenv("TEST_DB_RETRY_INTERVAL_MS"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				retryIntervalMs = n
			}
		}

		deadline := time.Now().Add(time.Duration(waitSeconds) * time.Second)
		var dbTemp *gorm.DB
		var err error
		for time.Now().Before(deadline) {
			dsnNoDB := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/?charset=utf8mb4&parseTime=True&loc=Local"
			dbTemp, err = gorm.Open(mysql.Open(dsnNoDB), &gorm.Config{})
			if err == nil {
				sqlDBTemp, _ := dbTemp.DB()
				if sqlDBTemp.Ping() == nil {
					// Create DB if not exists
					dbTemp.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
					dsn = dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
					break
				}
				sqlDBTemp.Close()
			}
			time.Sleep(time.Duration(retryIntervalMs) * time.Millisecond)
		}
		if dsn == "" {
			log.Fatalf("Failed to connect to DB server within %d seconds", waitSeconds)
		}
	} else {
		// fallback for local tests (use 127.0.0.1 to match CI service binding)
		dsnNoDB := "root:root@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
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
		dsn = "root:root@tcp(127.0.0.1:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
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
		// Drop the test database to keep environment clean
		if TestDBName != "" {
			// Execute drop using the current connection
			TestDB.Exec("DROP DATABASE IF EXISTS " + TestDBName)
		}
		sqlDB, _ := TestDB.DB()
		sqlDB.Close()
	}
}

func TruncateTables(db *gorm.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE TABLE palletized_products")
	db.Exec("TRUNCATE TABLE pallets")
	db.Exec("TRUNCATE TABLE pallet_racks")
	db.Exec("TRUNCATE TABLE products")
	db.Exec("TRUNCATE TABLE user_entities")
	db.Exec("TRUNCATE TABLE inventories")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}
