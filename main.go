package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/andreibugariu/BioShop-License/entity"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
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
	db.AutoMigrate(&entity.Orders_Products{})
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


	router.HandleFunc("/products", GetProducts).Methods("GET")//merge
	router.HandleFunc("/product/{id}", GetProduct).Methods("GET")//merge
	router.HandleFunc("/product", CreateProduct).Methods("POST")//merge
	router.HandleFunc("/product/{id}", DeleteProduct).Methods("DELETE")//merge
	router.HandleFunc("/product/{id}", UpdateProduct).Methods("PUT")//merge

	router.HandleFunc("/orders", GetOrders).Methods("GET")//merge
	router.HandleFunc("/order/{id}", GetOrder).Methods("GET")//merge
	router.HandleFunc("/order", CreateOrder).Methods("POST")//merge
	router.HandleFunc("/order/{id}", DeleteOrder).Methods("DELETE")//merge
	router.HandleFunc("/order/{id}", UpdateOrder).Methods("PUT")//merge

	router.HandleFunc("/orders_products", GetOrdersProducts).Methods("GET")//merge
	router.HandleFunc("/order_product/{id}", GetOrderProduct).Methods("GET")//merge
	router.HandleFunc("/order_product", CreateOrderProduct).Methods("POST")//merge
	router.HandleFunc("/order_product/{id}", DeleteOrderProduct).Methods("DELETE")//merge
	router.HandleFunc("/order_product/{id}", UpdateOrderProduct).Methods("PUT")//merge

	router.HandleFunc("/login", Login)//merge
	//router.HandleFunc("/home", Home).Methods("GET")
	router.HandleFunc("/logout", DeleteCookie)//merge

	log.Fatal(http.ListenAndServe(":8080", router))
}
///Login
var jwtKey = []byte("secret_key")


//Encode uses base64 as main encoding method.
func Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

//HasError uses logger to log the error if any has occured.
func HasError(err error, message string) bool {

	if err != nil {
		logrus.WithError(err).Error(message)
		return true
	}

	return false
}

type Claims struct {
	Email string `json:"username"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	reqbody := r.Body
	bodyBytes, err := ioutil.ReadAll(reqbody)

	if HasError(err,"Internal Error. Unable to read data") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}

	var user entity.User
	err = json.Unmarshal(bodyBytes, &user)

	if HasError(err,"Internal Error. Unmarshal problem") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}

	var userDB entity.User
	result := db.Find(&userDB, "email=?", user.Email)

	if result.RecordNotFound() {
		http.Error(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

	if user.Password != userDB.Password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Email: user.Email,
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

func IsAuth(w http.ResponseWriter, r *http.Request)  (bool, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return false, err
		}
		w.WriteHeader(http.StatusBadRequest)
		return false, err
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return false, err
		}
		w.WriteHeader(http.StatusBadRequest)
		return false, err
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false, err
	}
	
	return true, nil
}


func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
			Name:   "token",
			MaxAge: -1}
	http.SetCookie(w, &c)

	w.Write([]byte("old cookie deleted!\n"))
}

/////////////// API for USER
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
    
	isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var user []entity.User
	result := db.Find(&user)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&user)
    }

}


//Get specific user and their rentals
func GetUser(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
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
}

//Delete a specific user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
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
}

//Update a specific user by ID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
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
}

/////////////// API for FARMERS
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

    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var farmer []entity.Farmer
	result := db.Find(&farmer)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&farmer)
}

}


//Get specific farmer and their products
func GetFarmer(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
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
}

//Delete a specific user by ID
func DeleteFarmer(w http.ResponseWriter, r *http.Request) {
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
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
}

//Update a specific user by ID
func UpdateFarmer(w http.ResponseWriter, r *http.Request) {
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
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
}


/////////////// API for PRODUCT

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var product entity.Product

	json.NewDecoder(r.Body).Decode(&product)

	createProduct := db.Create(&product)
	err = createProduct.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&product)
	}
}
}

func GetProducts(w http.ResponseWriter, r *http.Request) {

    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var product []entity.Product
	result := db.Find(&product)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&product)
}

}



func GetProduct(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	params := mux.Vars(r)

	var product entity.Product
	var orders_products []entity.Orders_Products
	result := db.Where("id = ?", params["id"]).First(&product)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	db.Model(&product).Related(&orders_products)

	product.Orders_Products = orders_products

	json.NewEncoder(w).Encode(product)
}
}


func DeleteProduct(w http.ResponseWriter, r *http.Request) {
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	params := mux.Vars(r)

	var product entity.Product

	result := db.Where("id = ?", params["id"]).First(&product)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	result = db.Delete(&product)
	if result.Error != nil {
		http.Error(w, "can't delete product", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode("Product is succefully deleting")
}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var product entity.Product
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&product)

	result := db.Model(&entity.Product{}).Where("id= ?", params["id"]).Updates(product)
	if result.Error != nil {
		http.Error(w, "Can't update", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Product is succefully UPTDATE")
}
}


/////////////// API for ORDERS
func CreateOrder(w http.ResponseWriter, r *http.Request) {
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var order entity.Order

	json.NewDecoder(r.Body).Decode(&order)

	createOrder := db.Create(&order)
	err = createOrder.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&order)
	}
}
}


//Get all  order
func GetOrders(w http.ResponseWriter, r *http.Request) {

    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var order []entity.Order
	result := db.Find(&order)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&order)
}
	

}


//Get specific order
func GetOrder(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
    //Get order and the product
	isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	params := mux.Vars(r)

	var order entity.Order
	var orders_products []entity.Orders_Products
	result := db.Where("id = ?", params["id"]).First(&order)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	db.Model(&order).Related(&orders_products)

	order.Orders_Products = orders_products
	json.NewEncoder(w).Encode(order)
}
}

//Delete a specific order by ID
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	params := mux.Vars(r)

	var order entity.Order

	result := db.Where("id = ?", params["id"]).First(&order)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	result = db.Delete(&order)
	if result.Error != nil {
		http.Error(w, "can't delete order", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode("The order is succefully deleting")
}
}

//Update a specific order by ID
func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var order entity.Order
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&order)

	result := db.Model(&entity.Order{}).Where("id= ?", params["id"]).Updates(order)
	if result.Error != nil {
		http.Error(w, "Can't update", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Order is succefully UPTDATE")
}
}

/////////////// API for ORDERS_PRODUCT
func CreateOrderProduct(w http.ResponseWriter, r *http.Request) {
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var order_product entity.Orders_Products

	json.NewDecoder(r.Body).Decode(&order_product)

	createOrderProduct := db.Create(&order_product)
	err = createOrderProduct.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&order_product)
	}
}
}


//Get all  products
func GetOrdersProducts(w http.ResponseWriter, r *http.Request) {

    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var order_product []entity.Orders_Products
	result := db.Find(&order_product)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&order_product)
}

}


//Get specific user and their rentals
func GetOrderProduct(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	params := mux.Vars(r)

	var order_product entity.Orders_Products
	result := db.Where("id = ?", params["id"]).First(&order_product)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}
     
	json.NewEncoder(w).Encode(order_product)
}
}

//Delete a specific product by ID
func DeleteOrderProduct(w http.ResponseWriter, r *http.Request) {
    isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	params := mux.Vars(r)

	var order_product entity.Orders_Products

	result := db.Where("id = ?", params["id"]).First(&order_product)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	result = db.Delete(&order_product)
	if result.Error != nil {
		http.Error(w, "can't delete order", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode("The order_product is succefully deleting")
}
}

//Update a specific product by ID
func UpdateOrderProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err,"Error in authentication function") {
		http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
	var order_product entity.Orders_Products
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&order_product)

	result := db.Model(&entity.Orders_Products{}).Where("id= ?", params["id"]).Updates(order_product)
	if result.Error != nil {
		http.Error(w, "Can't update", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("succefully UPTDATE")
}
}

//var productsID []string
	//for _, num := range orders_products. {
	//    productsID = append(productsID,num)
	//  }
	/*var products []entity.Product

	    for _, productID := range id {
			var product entity.Product

			result := db.Where("id = ?", productID).First(&product)
		     if result.RecordNotFound() {
			  http.Error(w, "Not fount", http.StatusNotFound)
			  return
	  	   }
		 products := append(products, product)

	     }

		   ASA TREBUIE FACUT SI AICI

		ceva := []string{"Python", "Java", "C#", "Go", "Ruby"}
		var altceva[]string
		for _, a := range ceva{
		 var b string
		 b=a;
		 altceva =append(altceva,b)
		}
		fmt.Print(altceva)
	*/