package entity


type Product struct {

		ID  string   `gorm:"type:uuid;default:uuid_generate_v4();primary_key;column:id"`
		Name  string   `gorm:"type:varchar(255);not null"`
		Description  string   `gorm:"type:varchar(255);not null"`
		Price        string   `gorm:"type:varchar(255);not null"`
		Orders_Products   []Order_Product  

}
