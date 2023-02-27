package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	u "github.com/wuranxu/mouse/api/v1/util"
	"github.com/wuranxu/mouse/dao"
	"github.com/wuranxu/mouse/dto"
	"github.com/wuranxu/mouse/middleware"
	"github.com/wuranxu/mouse/model"
)

const (
	ErrInvalidUserNameOrPasswordCode = iota + 11000
	ErrGenerateTokenCode
)

var (
	ErrInvalidUserNameOrPassword = errors.New("invalid username/password")
	ErrGenerateToken             = errors.New("failed to create token")
)

type Api struct {
	app *gin.Engine
}

func New(app *gin.Engine) *Api {
	return &Api{app: app}
}

func (a *Api) AddRoute(middleware ...gin.HandlerFunc) {
	group := a.app.Group("/auth", middleware...)

	// route
	group.Handle("POST", "/login", u.Wrap(login))
}

func login(ctx *gin.Context) (any, error) {
	var user dto.LoginDto
	if err := ctx.ShouldBindJSON(&user); err != nil {
		return nil, err
	}
	resp := new(model.MouseUser)
	if err := dao.Conn.Debug().Find(&resp, `username = ? and password = ?`, user.Username, user.Password).Error; err != nil {
		return ErrInvalidUserNameOrPasswordCode, err
	}
	if resp.ID == 0 {
		// no user
		return ErrInvalidUserNameOrPasswordCode, ErrInvalidUserNameOrPassword
	}
	token, err := middleware.JWTUtil.CreateToken(middleware.CustomClaims{
		MouseUser: *resp,
	})
	if err != nil {
		return ErrGenerateTokenCode, ErrGenerateToken
	}
	resp.Token = token
	return resp, nil
}
