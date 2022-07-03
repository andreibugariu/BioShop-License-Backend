package api

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andreibugariu/BioShop-License/db"
	"github.com/andreibugariu/BioShop-License/entity"
	"github.com/andreibugariu/BioShop-License/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var err error

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

//Login user
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
	result := db.GetDB().Find(&userDB, "email=?", user.Email)

	if result.RecordNotFound() {
		http.Error(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

	err = util.CheckPassword(user.Password, userDB.Password)

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

//Check if it is authenticated
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

//Logout
func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(w, &c)

	w.Write([]byte("old cookie deleted!\n"))
}

//Get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {

	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		var user []entity.User
		result := db.GetDB().Find(&user)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(&user)
	}

}

//Get specific user by email
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
		result := db.GetDB().Where("email = ?", params["id"]).First(&user)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		result2 := db.GetDB().Where("user_email = ?", params["id"]).Find(&orders)
		if result2.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		user.Orders = orders

		json.NewEncoder(w).Encode(&user)
	}
}

//Create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user entity.User

	json.NewDecoder(r.Body).Decode(&user)

	var password string
	password, erro := util.HashPassword(user.Password)

	if erro != nil {
		json.NewEncoder(w).Encode(erro)
	}

	user.Password = password

	user.Category = "USER"
	createUser := db.GetDB().Create(&user)
	err = createUser.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&user.Password)
	}
}

//Delete user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	isAuth, err := IsAuth(w, r)
	if HasError(err, "Error in authentication function") {
		http.Error(w, "Internal Error. Please try again later", http.StatusInternalServerError)
		return
	}
	if isAuth {
		params := mux.Vars(r)

		var user entity.User

		result := db.GetDB().Where("id = ?", params["id"]).First(&user)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		result = db.GetDB().Delete(&user)
		if result.Error != nil {
			http.Error(w, "can't delete users", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode("Userul is succefully deleting")
	}
}

//Update user
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

		result := db.GetDB().Model(&entity.User{}).Where("id= ?", params["id"]).Updates(user)
		if result.Error != nil {
			http.Error(w, "Can't update", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode("Userul is succefully UPTDATE")
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
		result := db.GetDB().Where("email = ?", params["id"]).First(&user)
		if result.RecordNotFound() {
			http.Error(w, "Not fount", http.StatusNotFound)
			return
		}

		result2 := db.GetDB().Where("user_email = ? AND status = ?", params["id"], "active").Find(&orders)
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
			result := db.GetDB().Where("id = ?", a).First(&product)
			if result.RecordNotFound() {
				http.Error(w, "Not fount", http.StatusNotFound)
				return
			}

			full_products = append(full_products, product)
		}

		json.NewEncoder(w).Encode(full_products)
	}
}

//Get email from the cookie
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
