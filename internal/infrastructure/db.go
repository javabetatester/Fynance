package infrastructure

import (
	"log"

	"Fynance/internal/domain/goal"
	"Fynance/internal/domain/investment"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb() *gorm.DB {
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Falha ao conectar ao banco de dados: %v", err)
		log.Println("Certifique-se de que o PostgreSQL está rodando e as configurações estão corretas")
		log.Fatalf("Erro fatal na conexão com banco: %v", err)
	}

	log.Println("Conexão com banco de dados estabelecida com sucesso")

	runMigrations(db)

	return db
}

func runMigrations(db *gorm.DB) {
	log.Println("Executando migrations...")

	entities := []interface{}{
		&user.User{},
		&goal.Goal{},
		&transaction.Transaction{},
		&transaction.Category{},
		&investment.Investment{},
	}

	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			log.Printf("Erro ao migrar entidade %T: %v", entity, err)
			log.Fatalf("Falha na migração do banco de dados")
		}
	}

	log.Println("Migrations executadas com sucesso!")
}
