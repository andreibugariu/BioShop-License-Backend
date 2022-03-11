package main

import (
	"encoding/json"
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
	

	router.HandleFunc("/users", GetUsers).Methods("GET")//merge
	router.HandleFunc("/user/{id}", GetUser).Methods("GET")//merge
	router.HandleFunc("/user", CreateUser).Methods("POST")//merge
	router.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE")//merge
	router.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")//merge

	router.HandleFunc("/farmers", GetFarmers).Methods("GET")//merge
	router.HandleFunc("/farmer/{id}", GetFarmer).Methods("GET")//merge
	router.HandleFunc("/farmer", CreateFarmer).Methods("POST")//merge
	router.HandleFunc("/farmer/{id}", DeleteFarmer).Methods("DELETE")//merge
	router.HandleFunc("/farmer/{id}", UpdateFarmer).Methods("PUT")//merge


	log.Fatal(http.ListenAndServe(":8080", router))
}


//Create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user entity.User

	json.NewDecoder(r.Body).Decode(&user)

	createUser := db.Create(&user)
	err = createUser.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&user)
	}
}


//Get all  users
func GetUsers(w http.ResponseWriter, r *http.Request) {


	var user []entity.User
	result := db.Find(&user)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&user)

}


//Get specific user and their rentals
func GetUser(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.

	params := mux.Vars(r)

	var user entity.User
	var orders []entity.Order
	result := db.Where("id = ?", params["id"]).First(&user)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	db.Model(&user).Related(&orders)

	user.Orders = orders

	json.NewEncoder(w).Encode(user)
}

//Delete a specific user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var user entity.User

	result := db.Where("id = ?", params["id"]).First(&user)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	result = db.Delete(&user)
	if result.Error != nil {
		http.Error(w, "can't delete users", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode("Userul is succefully deleting")
}

//Update a specific user by ID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	
	var user entity.User
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&user)

	result := db.Model(&entity.User{}).Where("id= ?", params["id"]).Updates(user)
	if result.Error != nil {
		http.Error(w, "Can't update", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Userul is succefully UPTDATE")
}

//Create a new farmer
func CreateFarmer(w http.ResponseWriter, r *http.Request) {

	var farmer entity.Farmer

	json.NewDecoder(r.Body).Decode(&farmer)

	createFarmer := db.Create(&farmer)
	err = createFarmer.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&farmer)
	}
}


//Get all  farmers
func GetFarmers(w http.ResponseWriter, r *http.Request) {


	var farmer []entity.Farmer
	result := db.Find(&farmer)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&farmer)

}


//Get specific farmer and their products
func GetFarmer(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.

	params := mux.Vars(r)

	var farmer entity.Farmer
	var products []entity.Product
	result := db.Where("id = ?", params["id"]).First(&farmer)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	db.Model(&farmer).Related(&products)

	farmer.Products = products

	json.NewEncoder(w).Encode(farmer)
}

//Delete a specific user by ID
func DeleteFarmer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var farmer entity.Farmer

	result := db.Where("id = ?", params["id"]).First(&farmer)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	result = db.Delete(&farmer)
	if result.Error != nil {
		http.Error(w, "can't delete farmer", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode("Farmer is succefully deleting")
}

//Update a specific user by ID
func UpdateFarmer(w http.ResponseWriter, r *http.Request) {
	
	var farmer entity.Farmer
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&farmer)

	result := db.Model(&entity.Farmer{}).Where("id= ?", params["id"]).Updates(farmer)
	if result.Error != nil {
		http.Error(w, "Can't update", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("The Farmer is succefully UPTDATE")
}
