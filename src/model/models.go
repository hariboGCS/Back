package model

var authority bool

type User struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsTeacher  bool   `json:"is_teacher"`
	TotalScore int    `json:"total_score"`
	Token      string `json:"token"`
}

type Score struct {
	Email  string `json:"email"`
	Kind   int    `json:"kind"`
	Detail string `json:"detail"`
	Score  int    `json:"score"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}
