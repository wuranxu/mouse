package dto

type LoginDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Captcha  string `json:"captcha"`
}

type RegisterDto struct {
	Name          string `json:"name" binding:"required"`
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	PasswordAgain string `json:"passwordAgain" binding:"required,eqfield=Password"`
	Email         string `json:"email" binding:"required,email"`
}
