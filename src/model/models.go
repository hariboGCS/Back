package model

type Student struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Pass      string `json:"password"`
	LoggedIn  bool   `json:"loggedin"`
	IsTeacher bool   `json:"is_teacher"`
}
type Score struct {
	//엄청많음~~
}
