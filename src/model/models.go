package model

//User struct
type User struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsTeacher  bool   `json:"is_teacher"`
	TotalScore int    `json:"total_score"`
	Token      string `json:"token"`
}

//Score struct for User
type Score struct {
	Email  string `json:"email"`
	Kind   int    `json:"kind"`
	Detail string `json:"detail"`
	Score  int    `json:"score"`
}

//ResponseResult Mongo db connect
type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}
