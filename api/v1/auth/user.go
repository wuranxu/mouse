package auth

import (
	"crypto/md5"
	"encoding/hex"
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
	ErrParseParamCode
	ErrRegisterUserCode
	ErrUpdateLastLoginAtCode
	ErrUnAuthorizationCode
)

var (
	ErrInvalidUserNameOrPassword = errors.New("invalid username or password")
	ErrGenerateToken             = errors.New("failed to create token")
	ErrParseParam                = errors.New("wrong parameters")
	ErrRegisterUser              = errors.New("register error")
	ErrUpdateLastLoginAt         = errors.New("update loginAt error")
	ErrUnAuthorization           = errors.New("user not login")
)

const (
	salt       = "mouse-server"
	mouseToken = "mouse_token"
)

type Api struct {
	app *gin.Engine
}

func New(app *gin.Engine) *Api {
	return &Api{app: app}
}

func (a *Api) AddRoute(middleware ...gin.HandlerFunc) {
	group := a.app.Group("/auth", middleware...)

	group.Handle("POST", "/login", u.Wrap(login))
	group.Handle("POST", "/register", u.Wrap(register))
	group.Handle("GET", "/currentUser", u.Wrap(query))
}

func login(ctx *gin.Context) (any, error) {
	var user dto.LoginDto
	if err := ctx.ShouldBindJSON(&user); err != nil {
		return nil, u.ErrWrap(ErrParseParam, err)
	}
	resp := new(model.MouseUser)
	secret := md5.Sum([]byte(user.Password + salt))
	pwd := hex.EncodeToString(secret[:])
	if err := dao.Conn.MustFind(&resp, `username = ? and password = ?`, user.Username, pwd); err != nil {
		return ErrInvalidUserNameOrPasswordCode, ErrInvalidUserNameOrPassword
	}
	token, err := middleware.JWTUtil.CreateToken(middleware.CustomClaims{
		MouseUser: *resp,
	})
	if err != nil {
		return ErrGenerateTokenCode, ErrGenerateToken
	}
	resp.LastLoginAt = model.Now()
	if err = dao.Conn.Save(resp); err != nil {
		return ErrUpdateLastLoginAtCode, ErrUpdateLastLoginAt
	}
	ctx.SetCookie(mouseToken, token, 3600*8, "/", ctx.Request.Host, true, true)
	return resp, nil
}

func register(ctx *gin.Context) (any, error) {
	var data dto.RegisterDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		return ErrParseParamCode, u.ErrWrap(ErrParseParam, err)
	}
	secret := md5.Sum([]byte(data.Password + salt))
	pwd := hex.EncodeToString(secret[:])
	role := 0
	user := &model.MouseUser{
		Name:        data.Name,
		Username:    data.Username,
		Email:       data.Email,
		Password:    pwd,
		LastLoginAt: model.Now(),
		Role:        &role,
	}
	if err := dao.CreateUser(user); err != nil {
		return ErrRegisterUserCode, u.ErrWrap(ErrRegisterUser, err)
	}
	return user, nil
}

func query(ctx *gin.Context) (any, error) {
	token, err := ctx.Cookie(mouseToken)
	if err != nil {
		return ErrUnAuthorizationCode, ErrUnAuthorization
	}
	parseToken, err := middleware.JWTUtil.ParseToken(token)
	if err != nil {
		return ErrUnAuthorizationCode, ErrUnAuthorization
	}
	return parseToken.MouseUser, nil
}
