package rest

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/ajtwoddltka/GCS/src/dblayer"
	"github.com/ajtwoddltka/GCS/src/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type handlerInterface interface {
	GetScore(c echo.Context)
	ReceiveScore(c echo.Context)
	Login(c echo.Context) error
	Signup(c echo.Context) error
	updateUser(c echo.Context) error
	SignOut(c echo.Context)
	deleteUser(c echo.Context) error
	GetNotice(c echo.Context)
	GetComplaints(c echo.Context)
}

var (
	users = map[int]*model.User{}
	seq   = 1 //id에 들어갈 값
)

type (
	handler struct {
		DB *mgo.Session
	}
)

func newHandler() (*handler, error) { //핸들러 init
	return new(handler), nil
}

func (h *handler) Signup(c echo.Context) (err error) {
	// Bind
	u := &model.User{ID: bson.NewObjectId().Hex()}
	if err = c.Bind(u); err != nil {
		return err
	}
	// Validate
	if u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or password"}
	}
	collection, err := dblayer.GetDBCollection()
	collection.InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}
	defer collection.Database().Client().Disconnect(context.TODO())

	return c.JSON(http.StatusCreated, u)
}

func (h *handler) Signin(c echo.Context) (err error) {
	// Bind
	u := new(model.User)
	if err = c.Bind(u); err != nil {
		return
	}
	filter := bson.M{"username": u.Username, "password": u.Password}
	collection, err := dblayer.GetDBCollection()
	err = collection.FindOne(context.TODO(), filter).Decode(&u)
	if err != nil {
		return err
		// return &echo.HTTPError{Code: http.StatusUnauthorized,Message:"invalid email or password"}
	}
	defer collection.Database().Client().Disconnect(context.TODO())
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	u.Token, err = token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	u.Password = "" // Don't send password
	return c.JSON(http.StatusOK, u)
}

func (h *handler) updateUser(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Username = u.Username
	return c.JSON(http.StatusOK, users[id])
}

func (h *handler) deleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}
