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
	"github.com/andreibugariu/BioShop-License/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB
var err error

const DIALECT = "postgres"
const HOST = "localhost"
const DBPORT = "5432"
const USER = "postgres"
const NAME = "postgres"
const PASSWORD = "metalgreu98" 

func main() {
    testEncodeO()
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
	db.AutoMigrate(&entity.Users_Products{})
	db.AutoMigrate(&entity.Product{})

	router := mux.NewRouter()

	router.HandleFunc("/users", GetUsers).Methods("GET")          
	router.HandleFunc("/user/{id}", GetUser).Methods("GET")       
	router.HandleFunc("/user",CreateUser).Methods("POST")        
	router.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE") 
	router.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")    

	router.HandleFunc("/farmers", GetFarmers).Methods("GET")          
	router.HandleFunc("/farmer/{id}", GetFarmer).Methods("GET")       
	router.HandleFunc("/farmer", CreateFarmer).Methods("POST")        
	router.HandleFunc("/farmer/{id}", DeleteFarmer).Methods("DELETE") 
	router.HandleFunc("/farmer/{id}", UpdateFarmer).Methods("PUT")    

	router.HandleFunc("/farmer_products/{id}",GetFarmerProducts).Methods("GET") 

	router.HandleFunc("/farmer_orders/{id}",GetFarmerOrders).Methods("GET") 

	router.HandleFunc("/products", GetProducts).Methods("GET")                          
	router.HandleFunc("/product/{id}", GetProduct).Methods("GET")                       
	router.HandleFunc("/product", CreateProduct).Methods("POST")                        
	router.HandleFunc("/product/{id}", DeleteProduct).Methods("DELETE")                 
	router.HandleFunc("/product/{id}", UpdateProduct).Methods("PUT")                    
	router.HandleFunc("/product_by_category/{id}", GetProductByCategory).Methods("GET") 
	router.HandleFunc("/search_name/{id}", GetSearchProducts).Methods("GET")
	router.HandleFunc("/search_name_farmers_products/{email}/{name}", GetSearchProductsFarmers).Methods("GET")

	router.HandleFunc("/product_by_category_farmer/{email}/{category}", GetProductByCategoryFarmer).Methods("GET") 

	router.HandleFunc("/users_products", GetOrdersProducts).Methods("GET")               
	router.HandleFunc("/orders/{id}", GetOrderProduct).Methods("GET")                    
	router.HandleFunc("/user_product/{id}", GetCart).Methods("GET")                      
	router.HandleFunc("/user_product", CreateOrderProduct).Methods("POST")               
	router.HandleFunc("/delete_user_product/{id}", DeleteOrderProduct).Methods("DELETE") 
	router.HandleFunc("/increment/{id}", IncrementOrderProduct).Methods("PUT")           
	router.HandleFunc("/decrement/{id}", DecrementOrderProduct).Methods("PUT")
	router.HandleFunc("/checkout/{id}", OrderCheckout).Methods("PUT")

	
	

	router.HandleFunc("/login", Login)         
	router.HandleFunc("/login_farmer", LoginFarmer)         
	router.HandleFunc("/logout", DeleteCookie) 
	router.HandleFunc("/get_user", GetEmailCookie)

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

	if HasError(err, "Internal Error. Unable to read data") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}

	var user entity.User
	err = json.Unmarshal(bodyBytes, &user)

	if HasError(err, "Internal Error. Unmarshal problem") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
    
	var userDB entity.User
	result := db.Find(&userDB, "email=?", user.Email)

	if result.RecordNotFound() {
		http.Error(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

	err = util.CheckPassword(user.Password,userDB.Password)


	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 100)

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
	result := db.Find(&userDB, "email=?", farmer.Email)

	if result.RecordNotFound() {
		http.Error(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

		err = util.CheckPassword(farmer.Password,userDB.Password)

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

func IsAuth(w http.ResponseWriter, r *http.Request) (bool, error) {
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

func GetEmailCookie(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user, along with their
	// username given in the token
	// w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Email)))
	json.NewEncoder(w).Encode(&claims.Email)
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
    
	var password string
	password, erro := util.HashPassword(user.Password)

	if erro != nil {
		json.NewEncoder(w).Encode(erro)
	}

	user.Password=password

	user.Category = "USER"
	createUser := db.Create(&user)
	err = createUser.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&user.Password)
	}
}

//Get all  users
func GetUsers(w http.ResponseWriter, r *http.Request) {

	// isAuth, err := IsAuth(w, r)
	// if HasError(err,"Error in authentication function") {
	// 	http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
	// 	return
	// }
	// if isAuth {
	var user []entity.User
	result := db.Find(&user)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&user)
	//}

}

//Get specific user and their rentals
func GetUser(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var user entity.User
		var orders []entity.Users_Products
		result := db.Where("email = ?", params["id"]).First(&user)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		result2 := db.Where("user_email = ?", params["id"]).Find(&orders)
		if result2.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		user.Orders = orders

		json.NewEncoder(w).Encode(&user)
	}
}

//Delete a specific user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
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
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
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
	farmer.Category = "FARMER"

	json.NewDecoder(r.Body).Decode(&farmer)
	var password string
	password, erro := util.HashPassword(farmer.Password)

	if erro != nil {
		json.NewEncoder(w).Encode(erro)
	}
    farmer.Password=password
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
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
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
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var farmer entity.Farmer
		var products []entity.Product
		result := db.Where("email = ?", params["id"]).First(&farmer)
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
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
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
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
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
		json.NewEncoder(w).Encode("The Farmer is succefully UPDATE")
	}
}

/////////////// API for PRODUCT

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
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

	// isAuth, err := IsAuth(w, r)
	// if HasError(err,"Error in authentication function") {
	// 	http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
	// 	return
	// }
	// if isAuth {
	var product []entity.Product
	result := db.Find(&product)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&product)
	//}

}

func GetProductsByCategory(w http.ResponseWriter, r *http.Request) {

	// isAuth, err := IsAuth(w, r)
	// if HasError(err,"Error in authentication function") {
	// 	http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
	// 	return
	// }
	// if isAuth {
	params := mux.Vars(r)

	var products []entity.Users_Products
	result := db.Where("id = ?", params["id"]).Find(&products)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(products)

	//}

}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
	// isAuth, err := IsAuth(w, r)
	// if HasError(err,"Error in authentication function") {
	// 	http.Error(w,"Internal Error. Please try again later", http.StatusInternalServerError)
	// 	return
	// }
	// if isAuth {
	params := mux.Vars(r)

	var product entity.Product
	var users_products []entity.Users_Products
	result := db.Where("id = ?", params["id"]).First(&product)
	if result.RecordNotFound() {
		http.Error(w, "Not fount", http.StatusNotFound)
		return
	}

	db.Model(&product).Related(&users_products)

	product.Users_Products = users_products

	json.NewEncoder(w).Encode(product)
	//}
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
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
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

//#Get products by category own by farmer
func GetProductByCategoryFarmer( w http.ResponseWriter, r *http.Request){
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
		result := db.Where("farmer_email = ?", email).Find(&products)
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
		result := db.Find(&products)
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
func GetSearchProductsFarmers(w http.ResponseWriter, r *http.Request){
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
		result := db.Where("product_name LIKE ? AND farmer_email= ?","%"+name+"%",email).Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(&products)
	}
}

func GetSearchProducts(w http.ResponseWriter, r *http.Request) { ////succes
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
		result := db.Where("product_name LIKE ?", "%"+name+"%").Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(&products)
	}
}

/////////////// API for ORDERS_PRODUCT
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

		result := db.Where("id = ?", user_product.ProductID).First(&product)

		
        if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		 fproduct := float64(product.Quantity)
         
		if fproduct <= user_product.Quantity{
             	json.NewEncoder(w).Encode("exceeded capacity")
				return
		}

        new_quantity:=fproduct-user_product.Quantity

		result2 := db.Model(&product).Select("quantity").Updates(map[string]interface{}{"quantity": new_quantity})
		if result2.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}



		createUserProduct := db.Create(&user_product)
		err = createUserProduct.Error
		if err != nil {
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(&user_product)
		}
	}
}

//Get all  products
func GetOrdersProducts(w http.ResponseWriter, r *http.Request) {

	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var order_product []entity.Users_Products
		result := db.Find(&order_product)
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
		result := db.Where("user_email = ? AND status = ?", params["id"],"active").Find(&orders)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(orders)
	}
}

/////////////////////////////////////////GET PRODUCT BY EMAIL OF FARMER 
func GetFarmerProducts(w http.ResponseWriter, r *http.Request){
	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var products []entity.Product
		result := db.Where("farmer_email = ?", params["id"]).Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(products)
	}
}

/////////////////////////////////////////GET PRODUCT BY EMAIL OF FARMER 
func GetFarmerOrders(w http.ResponseWriter, r *http.Request){
	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var products []entity.Product
		result := db.Where("farmer_email = ?", params["id"]).Find(&products)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		
		var id_products []string
          
		 for _, s := range products {
         id_products= append(id_products, s.ID)
      }

	  var full_orders []entity.Users_Products
		
      for _, a := range id_products {
			var orders []entity.Users_Products
			result := db.Where("product_id = ?", a).Find(&orders)
			if result.RecordNotFound() {
				http.Error(w, "Not fount", http.StatusNotFound)
				return
			}

			full_orders = append(full_orders, orders...)
		}




		json.NewEncoder(w).Encode(full_orders)
	}
}



//Delete a specific product by ID
func DeleteOrderProduct(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var order_product entity.Users_Products

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
////////////////////////////////////////////////////////////////

//Update a specific product by ID
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

		result := db.Model(&entity.Users_Products{}).Where("id= ?", params["id"]).First(&order_product)
		if result.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		var new_quantity = order_product.Quantity + 1

		result2 := db.Model(&order_product).Select("quantity").Updates(map[string]interface{}{"quantity": new_quantity})
		if result2.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(order_product.Quantity)
	}
}

func DecrementOrderProduct(w http.ResponseWriter, r *http.Request) {//merge trebuie implementat in front
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var order_product entity.Users_Products
		params := mux.Vars(r)
		// json.NewDecoder(r.Body).Decode(&order_product)

		result := db.Model(&entity.Users_Products{}).Where("id= ?", params["id"]).First(&order_product)
		if result.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		var new_quantity = order_product.Quantity - 1

		result2 := db.Model(&order_product).Select("quantity").Updates(map[string]interface{}{"quantity": new_quantity})
		if result2.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(order_product.Quantity)
	}
}

func OrderCheckout(w http.ResponseWriter, r *http.Request) {///merge trebuie implementat in front-end

	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var orders []entity.Users_Products
		result := db.Where("user_email = ?", params["id"]).Find(&orders)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}
         var new_status = "canceled"

		result = db.Model(&orders).Select("status").Updates(map[string]interface{}{"status": new_status})

		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(orders)
	}

}

////Get all cart items from a user
func GetCart(w http.ResponseWriter, r *http.Request) {
	//Check the credentials provided in the request. Also check for errors at authentication.
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var user entity.User
		var orders []entity.Users_Products
		result := db.Where("email = ?", params["id"]).First(&user)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		result2 := db.Where("user_email = ? AND status = ?", params["id"] ,"active").Find(&orders)
		if result2.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		user.Orders = orders

		var products []string
		var full_products []entity.Product
		for _, a := range user.Orders {
			products = append(products, a.ProductID)
		}

		for _, a := range products {
			var product entity.Product
			result := db.Where("id = ?", a).First(&product)
			if result.RecordNotFound() {
				http.Error(w, "Not fount", http.StatusNotFound)
				return
			}

			full_products = append(full_products, product)
		}

		json.NewEncoder(w).Encode(full_products)
	}
}
func testEncodeO(){
	var parol string
    parol, err :=util.HashPassword("parolaa")

	if err != nil {
		fmt.Print("eror")
	}

	err = bcrypt.CompareHashAndPassword([]byte(parol),[]byte("ddparolaa"))

	if err != nil {
		fmt.Println("eror")
	}else {
		fmt.Print("merge")
	}

	
    err = util.CheckPassword("parolaa",parol)
    
	if err != nil {
		fmt.Println("eror")
	}else {
		fmt.Println("merge")
	}





} 