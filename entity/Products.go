package entity

type Product struct {
	ID             string  `gorm:"type:uuid;default:uuid_generate_v4();primary_key;column:id"`
	ProductName    string  `gorm:"type:varchar(255);not null"`
	Description    string  `gorm:"type:varchar(255);not null"`
	Price          float64 `gorm:"type:float;not null"`
	Photo          string  `gorm:"type:text;not null"`
	Category       string  `gorm:"type:varchar(255);not null"`
	FarmerID       string  `gorm:"type:uuid;column:farmer_id;not null" validate:"required"`
	Users_Products []Users_Products
}
