package model

var authority bool

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	LoggedIn   bool   `json:"loggedin"`
	IsTeacher  bool   `json:"is_teacher"`
	TotalScore int    `json:"total_score"`
	Token      string `json:"token"`
}

type Score struct {
	Kind     string `json:"kind"`
	Major    int    `json:"major_score"`
	Tenjob   int    `json:"tenjob_score"`
	Humanity int    `json:"humanity_score"`
	Foreign  int    `json:"foreign_score"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}
