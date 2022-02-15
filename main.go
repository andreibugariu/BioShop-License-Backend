package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andreibugariu/BioShop-License/entity"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

/*var (
	user = &entity.User{

		FirstName:   "Andrei",
		LastName:    "Bugariu",
		Email:       "andreibugariu@gmail.com",
		Password:    "DJSDNASDKNA",
		Age:         21,
		Address:     "Timisoara",
		PhoneNumber: 90712987,
		BankDetails: "IBAN1293Y2138Y",
	}
)*/

var db *gorm.DB
var err error

const DIALECT = "postgres"
const HOST = "localhost"
const DBPORT = "5432"
const USER = "postgres"
const NAME = "postgres"
const PASSWORD = "metalgreu98" ///here introduce your password !!!!

func main() {

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", HOST, USER, NAME, PASSWORD, DBPORT)

	db, err = gorm.Open(DIALECT, dbURI)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Successfully connected to database")
	}

	defer db.Close()
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Farmer{})
	db.AutoMigrate(&entity.Order_Product{})
	db.AutoMigrate(&entity.Product{})
	db.AutoMigrate(&entity.Order{})

	/*createUser := db.Create(user)
	err = createUser.Error
	if err != nil {
		fmt.Println("are eroare")
	} else {
		fmt.Println("nu are  eroare")
	}*/

	router := mux.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

}
