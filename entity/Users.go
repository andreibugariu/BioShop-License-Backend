package entity

type User struct {
	ID          string           `gorm:"type:uuid;default:uuid_generate_v4();primary_key;column:id"`
	FirstName   string           `gorm:"type:varchar(255);not null"`
	LastName    string           `gorm:"type:varchar(255);not null"`
	Age         int64            `gorm:"type:int;not null"`
	Address     string           `gorm:"type:varchar(255);not null"`
	Email       string           `gorm:"type:varchar(255);not null;unique" validate:"required,email"`
	Password    string           `gorm:"type:varchar(255);not null" validate:"required,min=6"`
	PhoneNumber string           `gorm:"type:varchar(255);not null"`
	Category    string           `gorm:"type:varchar(255);not null"`
	Orders      []Users_Products `gorm:"foreignKey:UserProductID"`
}
