package testsetup

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// DBConfig holds the database connection configuration for tests
// You can extend this struct if needed for other setups
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// GetTestDBConfig reads environment variables and returns DBConfig
func GetTestDBConfig() DBConfig {
	return DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

// ConnectWithRetry tries to connect to the DB with retries
func ConnectWithRetry(cfg DBConfig, maxRetries int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		time.Sleep(delay)
	}
	return nil, fmt.Errorf("could not connect to DB after %d retries: %w", maxRetries, err)
}

// SeedAdminUser insere o usuário admin na tabela user_entities para testes e2e
func SeedAdminUser(db *sql.DB, name, email, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	_, err = db.Exec(`INSERT INTO user_entities (name, email, password, role) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE role=VALUES(role), password=VALUES(password)`,
		name, email, string(hash), role)
	if err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}
	return nil
}
