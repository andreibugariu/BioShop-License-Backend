package db

import (
	"fmt"

	"github.com/andreibugariu/BioShop-License/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB //database

const DIALECT = "postgres"
const HOST = "localhost"
const DBPORT = "5432"
const USER = "postgres"
const NAME = "postgres"
const PASSWORD = "metalgreu98"

func init() {
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", HOST, USER, NAME, PASSWORD, DBPORT) //Build connection string
	fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&entity.User{}, &entity.Farmer{}, &entity.Users_Products{}, &entity.Product{})
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
