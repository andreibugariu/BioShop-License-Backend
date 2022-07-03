package api

import (
	"encoding/json"
	"net/http"

	"github.com/andreibugariu/BioShop-License/db"
	"github.com/andreibugariu/BioShop-License/entity"
	"github.com/gorilla/mux"
)

//Create an order
func CreateOrderProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var user_product entity.Users_Products

		json.NewDecoder(r.Body).Decode(&user_product)
		user_product.Status = "active"

		var product entity.Product

		result := db.GetDB().Where("id = ?", user_product.ProductID).First(&product)

		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		fproduct := float64(product.Quantity)

		if fproduct <= user_product.Quantity {
			json.NewEncoder(w).Encode("exceeded capacity")
			return
		}

		new_quantity := fproduct - user_product.Quantity

		result2 := db.GetDB().Model(&product).Select("quantity").Updates(map[string]interface{}{"quantity": new_quantity})
		if result2.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}

		createUserProduct := db.GetDB().Create(&user_product)
		err = createUserProduct.Error
		if err != nil {
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(&user_product)
		}
	}
}

func GetOrdersProducts(w http.ResponseWriter, r *http.Request) {

	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var order_product []entity.Users_Products
		result := db.GetDB().Find(&order_product)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(&order_product)
	}

}

func GetOrderProduct(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var orders []entity.Users_Products
		result := db.GetDB().Where("user_email = ? AND status = ?", params["id"], "active").Find(&orders)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(orders)
	}
}

//Delete order
func DeleteOrderProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var order_product entity.Users_Products

		result := db.GetDB().Where("id = ?", params["id"]).First(&order_product)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		result = db.GetDB().Delete(&order_product)
		if result.Error != nil {
			http.Error(w, "can't delete order", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode("The order_product is succefully deleting")
	}
}

//Increases the quantity of the product in the order
func IncrementOrderProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var order_product entity.Users_Products
		params := mux.Vars(r)
		// json.NewDecoder(r.Body).Decode(&order_product)

		result := db.GetDB().Model(&entity.Users_Products{}).Where("id= ?", params["id"]).First(&order_product)
		if result.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		var new_quantity = order_product.Quantity + 1

		result2 := db.GetDB().Model(&order_product).Select("quantity").Updates(map[string]interface{}{"quantity": new_quantity})
		if result2.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(order_product.Quantity)
	}
}

//Decreases the quantity of the product in the order
func DecrementOrderProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var order_product entity.Users_Products
		params := mux.Vars(r)
		// json.NewDecoder(r.Body).Decode(&order_product)

		result := db.GetDB().Model(&entity.Users_Products{}).Where("id= ?", params["id"]).First(&order_product)
		if result.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		var new_quantity = order_product.Quantity - 1

		result2 := db.GetDB().Model(&order_product).Select("quantity").Updates(map[string]interface{}{"quantity": new_quantity})
		if result2.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(order_product.Quantity)
	}
}

//complete the order
func OrderCheckout(w http.ResponseWriter, r *http.Request) {

	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var orders []entity.Users_Products
		result := db.GetDB().Where("user_email = ?", params["id"]).Find(&orders)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}
		var new_status = "canceled"

		result = db.GetDB().Model(&orders).Select("status").Updates(map[string]interface{}{"status": new_status})

		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(orders)
	}

}
