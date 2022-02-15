package entity


type Order struct{

	ID          string   `gorm:"type:uuid;default:uuid_generate_v4();primary_key;column:id"`
	Year_of_release int64   `gorm:"type:varchar(255);not null"`
    Orders_Products   []Order_Product  

}