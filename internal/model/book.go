package model

type Book struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
	Stock string `json:"stock"`
}

type ParseBook struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `type:text;not null"`
	Price string `type:text;not null"`
	Stock string `type:text;not null"`
}
