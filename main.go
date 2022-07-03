package main

import (
	"log"
	"net/http"

	"github.com/andreibugariu/BioShop-License/api"

	"github.com/gorilla/mux"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/users", api.GetUsers).Methods("GET")
	router.HandleFunc("/user/{id}", api.GetUser).Methods("GET")
	router.HandleFunc("/user", api.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", api.DeleteUser).Methods("DELETE")
	router.HandleFunc("/user/{id}", api.UpdateUser).Methods("PUT")
	router.HandleFunc("/login", api.Login)

	router.HandleFunc("/farmers", api.GetFarmers).Methods("GET")
	router.HandleFunc("/farmer/{id}", api.GetFarmer).Methods("GET")
	router.HandleFunc("/farmer", api.CreateFarmer).Methods("POST")
	router.HandleFunc("/farmer/{id}", api.DeleteFarmer).Methods("DELETE")
	router.HandleFunc("/farmer/{id}", api.UpdateFarmer).Methods("PUT")

	router.HandleFunc("/farmer_products/{id}", api.GetFarmerProducts).Methods("GET")

	router.HandleFunc("/farmer_orders/{id}", api.GetFarmerOrders).Methods("GET")

	router.HandleFunc("/products", api.GetProducts).Methods("GET")
	router.HandleFunc("/product/{id}", api.GetProduct).Methods("GET")
	router.HandleFunc("/product", api.CreateProduct).Methods("POST")
	router.HandleFunc("/product/{id}", api.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/product/{id}", api.UpdateProduct).Methods("PUT")
	router.HandleFunc("/product_by_category/{id}", api.GetProductByCategory).Methods("GET")
	router.HandleFunc("/search_name/{id}", api.GetSearchProducts).Methods("GET")
	router.HandleFunc("/search_name_farmers_products/{email}/{name}", api.GetSearchProductsFarmers).Methods("GET")

	router.HandleFunc("/product_by_category_farmer/{email}/{category}", api.GetProductByCategoryFarmer).Methods("GET")

	router.HandleFunc("/users_products", api.GetOrdersProducts).Methods("GET")
	router.HandleFunc("/orders/{id}", api.GetOrderProduct).Methods("GET")
	router.HandleFunc("/user_product/{id}", api.GetCart).Methods("GET")
	router.HandleFunc("/user_product", api.CreateOrderProduct).Methods("POST")
	router.HandleFunc("/delete_user_product/{id}", api.DeleteOrderProduct).Methods("DELETE")
	router.HandleFunc("/increment/{id}", api.IncrementOrderProduct).Methods("PUT")
	router.HandleFunc("/decrement/{id}", api.DecrementOrderProduct).Methods("PUT")
	router.HandleFunc("/checkout/{id}", api.OrderCheckout).Methods("PUT")

	router.HandleFunc("/login_farmer", api.LoginFarmer)
	router.HandleFunc("/logout", api.DeleteCookie)
	router.HandleFunc("/get_user", api.GetEmailCookie)

	log.Fatal(http.ListenAndServe(":8080", router))

}
