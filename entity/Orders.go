package entity

type Order struct {
	ID              string `gorm:"type:uuid;default:uuid_generate_v4();primary_key;column:id"`
	UserID          string `gorm:"type:uuid;column:user_id;not null" validate:"required"`
	Price           int64  `gorm:"type:varchar(255);not null"`
	Orders_Products []Orders_Products
}
