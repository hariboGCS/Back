package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ajtwoddltka/GCS/src/dblayer"
	"github.com/ajtwoddltka/GCS/src/model"
	"github.com/unrolled/render"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var rd *render.Render

type HandlerInterface interface {
	GetMainPage(w http.ResponseWriter, r *http.Request)
	GetScore(w http.ResponseWriter, r *http.Request)
	GetRankPage(w http.ResponseWriter, r *http.Request)
	GetProfile(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
	GetStandard(w http.ResponseWriter, r *http.Request)
	GetNotice(w http.ResponseWriter, r *http.Request)
	GetComplaints(w http.ResponseWriter, r *http.Request)
}
type Handler struct {
	handler HandlerInterface
}

func NewHandler() (*Handler, error) {
	return new(Handler), nil
}
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var result model.User
	var res model.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username = claims["username"].(string)
		result.FirstName = claims["firstname"].(string)
		result.LastName = claims["lastname"].(string)

		json.NewEncoder(w).Encode(result)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}

func (h *Handler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Main Page")
}
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	collection, err := dblayer.GetDBCollection()

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

			if err != nil {
				res.Error = "Error While Hashing Password, Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}
			user.Password = string(hash)

			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Error = "Error While Creating User, Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}
			res.Result = "Registration Successful"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Result = "Username already Exists!!"
	json.NewEncoder(w).Encode(res)
	return
}
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	collection, err := db.GetDBCollection()

	if err != nil {
		log.Fatal(err)
	}
	var result model.User
	var res model.ResponseResult

	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

	if err != nil {
		res.Error = "Invalid username"
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		res.Error = "Invalid password"
		json.NewEncoder(w).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  result.Username,
		"firstname": result.FirstName,
		"lastname":  result.LastName,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		res.Error = "Error while generating token,Try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""

	json.NewEncoder(w).Encode(result)
}
func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) GetStandard(w http.ResponseWriter, r *http.Request) {

}
