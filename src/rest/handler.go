package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ajtwoddltka/GCS/src/dblayer"
	"github.com/ajtwoddltka/GCS/src/model"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type HandlerInterface interface {
	GetMainPage(w http.ResponseWriter, r *http.Request)
	GetScore(w http.ResponseWriter, r *http.Request)
	ReceiveScore(w http.ResponseWriter, r *http.Request)
	GetProfile(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
	GetNotice(w http.ResponseWriter, r *http.Request)
	GetComplaints(w http.ResponseWriter, r *http.Request)
}
type Handler struct {
	handler HandlerInterface //핸들로 인터페이스 structure 정의
}

func NewHandler() (*Handler, error) {
	return new(Handler), nil
}

func (h *Handler) GetMainPage(w http.ResponseWriter, r *http.Request) { //mainpage
	fmt.Fprintln(w, "Main Page")
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	user := &model.User{"", "", "", false, 0, ""}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	score := &model.Score{"", 0, "", 0}
	body, _ = ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &score)

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(score)
		return
	}
	score.Email = user.Email

	collection, err := dblayer.GetDBCollection()

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			fmt.Println(user.Password)
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
		}
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
	} else {
		res.Result = "Email	already Exists!!"
		json.NewEncoder(w).Encode(res)
		return
	}
	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Error = "Error While Creating User, Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}
			res.Result = "Insert Score to users"
			json.NewEncoder(w).Encode(res)
			return
		}
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)

	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := &model.User{"", "", "", false, 0, ""}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	collection, err := dblayer.GetDBCollection()

	if err != nil {
		log.Fatal(err)
	}
	result := &model.User{"", "", "", true, 0, ""}
	var res model.ResponseResult

	err = collection.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}}).Decode(&result)

	if err != nil {
		res.Result = "Invalid email"
		json.NewEncoder(w).Encode(res)
		http.Redirect(w, r, "http://localhost:3000/login", http.StatusTemporaryRedirect)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		res.Result = "Invalid password"
		json.NewEncoder(w).Encode(res)
		http.Redirect(w, r, "http://localhost:3000/login", http.StatusTemporaryRedirect)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": result.Email,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		res.Error = "Error while generating token,Try again"
		json.NewEncoder(w).Encode(res)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	email := user.Email

	rs, err := collection.UpdateOne(
		ctx,
		bson.M{"email": email},
		bson.D{
			{"$set", bson.D{{Key: "loggedin", Value: true}}},
		},
	)
	if err != nil {
		fmt.Fprintf(w, "db에러")
	}
	fmt.Printf("Updated %v loggedin\n", rs.ModifiedCount)

	result.Token = tokenString
	result.Password = ""

	json.NewEncoder(w).Encode(result)
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

		json.NewEncoder(w).Encode(result)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
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
		log.Fatal(err)
	}
	var result model.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	email := user.Email
	rs, err := collection.UpdateOne(
		ctx,
		bson.M{"email": email},
		bson.D{
			{"$set", bson.D{primitive.E{Key: "loggedin", Value: false}}},
		},
	)
	fmt.Printf("Updated %v signout\n", rs.ModifiedCount)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": result.Username,
	})
	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		res.Error = "Error while generating token,Try again"
		json.NewEncoder(w).Encode(res)
		return

	}
	result.Token = tokenString
}

func (h *Handler) ReceiveScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	score := &model.Score{"", 0, "", 0}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &score)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(score)
		return
	}
	collection, err := dblayer.GetDBCollection()

	var user model.Score

	err = collection.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}}).Decode(&score)

	if score.Email != user.Email {
		res.Error = "보내준 email 값이 다름"
		json.NewEncoder(w).Encode(res)
		return
	}

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	if err != nil {
		res.Error = "찾을 수 없습니다"
		json.NewEncoder(w).Encode(res)
		return
	}

	_, err = collection.InsertOne(context.TODO(), score)
}

func (h *Handler) GetScore(w http.ResponseWriter, r *http.Request) {
	user := &model.User{"", "", "", true, 0, ""}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	var res model.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var result model.Score

	collection, err := dblayer.GetDBCollection()
	err = collection.FindOne(context.TODO(), bson.D{{Key: "email", Value: result.Email}}).Decode(&result)

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	_, err = collection.InsertOne(context.TODO(), result)
}
