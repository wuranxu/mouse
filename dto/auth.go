package dto

type LoginDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Captcha  string `json:"captcha"`
}

type RegisterDto struct {
	Name          string `json:"name" validate:"required"`
	Username      string `json:"username" validate:"required"`
	Password      string `json:"password" validate:"required"`
	PasswordAgain string `json:"passwordAgain" validate:"required,eqfield=Password"`
	Email         string `json:"email" validate:"required,email"`
}
