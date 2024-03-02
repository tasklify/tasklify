package database

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func (s *database) CreateProduct(email string, password string) error {

	// Create
	s.Create(&Product{Code: "D42", Price: 100})

	return nil
}
