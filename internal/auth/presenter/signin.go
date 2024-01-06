package presenter

type LoginInput struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}