package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andreibugariu/BioShop-License/db"
	"github.com/andreibugariu/BioShop-License/entity"
	"github.com/andreibugariu/BioShop-License/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

//Login farmer
func LoginFarmer(w http.ResponseWriter, r *http.Request) {
	reqbody := r.Body
	bodyBytes, err := ioutil.ReadAll(reqbody)

	if HasError(err, "Internal Error. Unable to read data") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}

	var farmer entity.Farmer
	err = json.Unmarshal(bodyBytes, &farmer)

	if HasError(err, "Internal Error. Unmarshal problem") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}

	var userDB entity.Farmer
	result := db.GetDB().Find(&userDB, "email=?", farmer.Email)

	if result.RecordNotFound() {
		http.Error(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

	err = util.CheckPassword(farmer.Password, userDB.Password)

	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 100)

	claims := &Claims{
		Email: farmer.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

}

func CreateFarmer(w http.ResponseWriter, r *http.Request) {

	var farmer entity.Farmer
	farmer.Category = "FARMER"

	json.NewDecoder(r.Body).Decode(&farmer)
	var password string
	password, erro := util.HashPassword(farmer.Password)

	if erro != nil {
		json.NewEncoder(w).Encode(erro)
	}
	farmer.Password = password
	createFarmer := db.GetDB().Create(&farmer)
	err = createFarmer.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&farmer)
	}
}

//Update a specific user by ID
func UpdateFarmer(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var farmer entity.Farmer
		params := mux.Vars(r)
		json.NewDecoder(r.Body).Decode(&farmer)

		result := db.GetDB().Model(&entity.Farmer{}).Where("id= ?", params["id"]).Updates(farmer)
		if result.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode("The Farmer is succefully UPDATE")
	}
}

//Get all farmers
func GetFarmers(w http.ResponseWriter, r *http.Request) {

	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var farmer []entity.Farmer
		result := db.GetDB().Find(&farmer)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(&farmer)
	}

}

//Get farmer by email
func GetFarmer(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var farmer entity.Farmer
		var products []entity.Product
		result := db.GetDB().Where("email = ?", params["id"]).First(&farmer)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		db.GetDB().Model(&farmer).Related(&products)

		farmer.Products = products

		json.NewEncoder(w).Encode(farmer)
	}
}

//DeleteFarmer
func DeleteFarmer(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var farmer entity.Farmer

		result := db.GetDB().Where("id = ?", params["id"]).First(&farmer)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		result = db.GetDB().Delete(&farmer)
		if result.Error != nil {
			http.Error(w, "can't delete farmer", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode("Farmer is succefully deleting")
	}
}

//Get farmers products
func GetFarmerProducts(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var products []entity.Product
		result := db.GetDB().Where("farmer_email = ?", params["id"]).Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(products)
	}
}

func GetFarmerOrders(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var products []entity.Product
		result := db.GetDB().Where("farmer_email = ?", params["id"]).Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		var id_products []string

		for _, s := range products {
			id_products = append(id_products, s.ID)
		}

		var full_orders []entity.Users_Products

		for _, a := range id_products {
			var orders []entity.Users_Products
			result := db.GetDB().Where("product_id = ?", a).Find(&orders)
			if result.RecordNotFound() {
				http.Error(w, "Not fount", http.StatusNotFound)
				return
			}

			full_orders = append(full_orders, orders...)
		}

		json.NewEncoder(w).Encode(full_orders)
	}
}
