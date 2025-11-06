package infrastructure

import (
	"Fynance/internal/config"
	"Fynance/internal/domain/goal"
	"Fynance/internal/domain/investment"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	"Fynance/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		logger.Error().
			Err(err).
			Str("host", cfg.Database.Host).
			Int("port", cfg.Database.Port).
			Str("database", cfg.Database.DBName).
			Msg("Falha ao conectar ao banco de dados")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error().Err(err).Msg("Falha ao obter instância do banco de dados")
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	logger.Info().
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("database", cfg.Database.DBName).
		Msg("Conexão com banco de dados estabelecida com sucesso")

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *gorm.DB) error {
	logger.Info().Msg("Executando migrations...")

	entities := []interface{}{
		&user.User{},
		&goal.Goal{},
		&transaction.Transaction{},
		&transaction.Category{},
		&investment.Investment{},
	}

	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			logger.Error().
				Err(err).
				Str("entity", getEntityName(entity)).
				Msg("Erro ao migrar entidade")
			return err
		}
	}

	logger.Info().Msg("Migrations executadas com sucesso!")
	return nil
}

func getEntityName(entity interface{}) string {
	switch entity.(type) {
	case *user.User:
		return "User"
	case *goal.Goal:
		return "Goal"
	case *transaction.Transaction:
		return "Transaction"
	case *transaction.Category:
		return "Category"
	case *investment.Investment:
		return "Investment"
	default:
		return "Unknown"
	}
}
