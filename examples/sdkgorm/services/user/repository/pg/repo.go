package pg

import "gorm.io/gorm"

type pgRepo struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *pgRepo {
	return &pgRepo{db: db}
}
