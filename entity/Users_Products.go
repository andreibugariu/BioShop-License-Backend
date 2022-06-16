package entity

type Users_Products struct {
	ID        string  `gorm:"type:uuid;default:uuid_generate_v4();primary_key;column:id"`
	UserEmail string  `gorm:"type:varchar(255);not null"`
	ProductID string  `gorm:"type:uuid;column:product_id;not null" validate:"required"`
	ProductName string  `gorm:"type:varchar(255);not null"`
	ProductPrice  float64 `gorm:"type:float;not null"`
	Quantity  float64 `gorm:"type:float;not null"`
    Status string  `gorm:"type:varchar(255);not null"`
}
