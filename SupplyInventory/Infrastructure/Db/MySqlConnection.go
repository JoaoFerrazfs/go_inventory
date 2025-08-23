package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := "root:root@tcp(db:3306)/inventory?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar no banco:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	// Testa a conexão
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Erro ao pingar o banco:", err)
	}

	log.Println("Conexão com o banco realizada com sucesso")
	return db
}
