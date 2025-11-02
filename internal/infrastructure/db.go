package infrastructure

import (
	"Fynance/internal/domain/goal"
	"Fynance/internal/domain/investment"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb() *gorm.DB {
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("fail to connect to database")
	}

	err = db.AutoMigrate(
		&user.User{},
		&transaction.Transaction{},
		&transaction.Category{},
		&investment.Investment{},
		&goal.Goal{},
	)
	if err != nil {
		log.Printf("Erro no AutoMigrate: %v", err)
		panic("fail to migrate database")
	}

	log.Println("AutoMigrate executado com sucesso!")
	return db
}
