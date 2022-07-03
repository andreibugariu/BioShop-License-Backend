package api

import (
	"encoding/json"
	"net/http"

	"github.com/andreibugariu/BioShop-License/db"
	"github.com/andreibugariu/BioShop-License/entity"
	"github.com/gorilla/mux"
)

//Create product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var product entity.Product

		json.NewDecoder(r.Body).Decode(&product)

		createProduct := db.GetDB().Create(&product)
		err = createProduct.Error
		if err != nil {
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(&product)
		}
	}
}

//Get all products
func GetProducts(w http.ResponseWriter, r *http.Request) {

	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var product []entity.Product
		result := db.GetDB().Find(&product)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(&product)
	}

}

//Get products by category
func GetProductsByCategory(w http.ResponseWriter, r *http.Request) {

	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var products []entity.Users_Products
		result := db.GetDB().Where("id = ?", params["id"]).Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(products)

	}

}

//Get product by id
func GetProduct(w http.ResponseWriter, r *http.Request) {
	// Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var product entity.Product
		var users_products []entity.Users_Products
		result := db.GetDB().Where("id = ?", params["id"]).First(&product)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		db.GetDB().Model(&product).Related(&users_products)

		product.Users_Products = users_products

		json.NewEncoder(w).Encode(product)
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var product entity.Product

		result := db.GetDB().Where("id = ?", params["id"]).First(&product)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		result = db.GetDB().Delete(&product)
		if result.Error != nil {
			http.Error(w, "can't delete product", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode("Product is succefully deleting")
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var product entity.Product
		params := mux.Vars(r)
		json.NewDecoder(r.Body).Decode(&product)

		result := db.GetDB().Model(&entity.Product{}).Where("id= ?", params["id"]).Updates(product)
		if result.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode("Product is succefully UPTDATE")
	}
}

//Get products by category own by farmer
func GetProductByCategoryFarmer(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		email := params["email"]
		category := params["category"]
		var products []entity.Product
		result := db.GetDB().Where("farmer_email = ?", email).Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}
		var category_products []entity.Product
		for _, product := range products {
			if product.Category == category {
				category_products = append(category_products, product)
			}

		}
		if len(category_products) == 0 {
			http.Error(w, "Wrong category", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(&category_products)
	}
}

func GetProductByCategory(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.

	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)
		var products []entity.Product
		result := db.GetDB().Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}
		var category_products []entity.Product
		for _, product := range products {
			if product.Category == params["id"] {
				category_products = append(category_products, product)
			}

		}
		if len(category_products) == 0 {
			http.Error(w, "Wrong category", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(&category_products)
	}
}

//Get search produc own by farmer
func GetSearchProductsFarmers(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)
		email := params["email"]
		name := params["name"]
		var products []entity.Product
		result := db.GetDB().Where("product_name LIKE ? AND farmer_email= ?", "%"+name+"%", email).Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(&products)
	}
}

func GetSearchProducts(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)
		var products []entity.Product
		var name = params["id"]
		result := db.GetDB().Where("product_name LIKE ?", "%"+name+"%").Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(&products)
	}
}
