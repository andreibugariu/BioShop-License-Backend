package entity


type Orders_Products struct {

	ID       string `gorm:"type:uuid;default:uuid_generate_v4();primary_key;column:id"`
	OrderID string `gorm:"type:uuid;column:order_id;not null" validate:"required"`
	ProductID  string `gorm:"type:uuid;column:product_id;not null" validate:"required"`


}