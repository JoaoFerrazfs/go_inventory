package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "root:root@tcp(db:3306)/inventory?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar no banco:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal(err)
	}

	// Testa a conexão
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Erro ao pingar o banco:", err)
	}

	log.Println("Conexão com o banco realizada com sucesso")
}
