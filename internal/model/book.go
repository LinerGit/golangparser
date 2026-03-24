package model

type Book struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Price string `json:"price"`
	Stock string `json:"stock"`
}
