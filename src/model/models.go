package model

import "time"

var authority bool

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Pass      string    `json:"password"`
	LoggedIn  bool      `json:"loggedin"`
	IsTeacher bool      `json:"is_teacher"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}
