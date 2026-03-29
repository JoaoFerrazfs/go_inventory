package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	var dsn string

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")

	fmt.Println("Conectando ao banco de dados MySQL...", dbHost, dbPort, dbName, dbUser)

	// Se todas as variáveis estiverem definidas, monta o DSN
	if dbHost != "" && dbPort != "" && dbName != "" && dbUser != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbPass, dbHost, dbPort, dbName)
	} else {
		// fallback local padrão
		dsn = "root:root@tcp(db:3306)/inventory?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar no banco:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Erro ao pingar o banco:", err)
	}

	log.Println("Conexão com o banco realizada com sucesso")
	return db
}
