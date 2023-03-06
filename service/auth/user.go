package auth

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/conf"
	"github.com/wuranxu/mouse/dao"
	"github.com/wuranxu/mouse/dto"
	"github.com/wuranxu/mouse/exception"
	"github.com/wuranxu/mouse/middleware"
	"github.com/wuranxu/mouse/model"
	"github.com/wuranxu/mouse/utils/request"
)

const (
	ErrInvalidUserNameOrPasswordCode = iota + 11000
	ErrGenerateTokenCode
	ErrRegisterUserCode
	ErrUpdateLastLoginAtCode
	ErrUnAuthorizationCode
)

var (
	ErrInvalidUserNameOrPassword = exception.Err("invalid username or password")
	ErrGenerateToken             = exception.Err("failed to create token")
	ErrRegisterUser              = exception.Err("register error")
	ErrUpdateLastLoginAt         = exception.Err("update loginAt error")
	ErrUnAuthorization           = exception.Err("user not login")
)

const (
	salt = "mouse-server"
)

func Login(ctx *gin.Context) (any, error) {
	user := request.GetJson[dto.LoginDto](ctx)
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
	ctx.SetCookie(conf.MouseToken, token, 3600*8, "/", ctx.Request.Host, true, true)
	return resp, nil
}

func Register(ctx *gin.Context) (any, error) {
	data := request.GetJson[dto.RegisterDto](ctx)
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
		return ErrRegisterUserCode, ErrRegisterUser.New(err)
	}
	return user, nil
}

func Query(ctx *gin.Context) (any, error) {
	token, err := ctx.Cookie(conf.MouseToken)
	if err != nil {
		return ErrUnAuthorizationCode, ErrUnAuthorization
	}
	parseToken, err := middleware.JWTUtil.ParseToken(token)
	if err != nil {
		return ErrUnAuthorizationCode, ErrUnAuthorization
	}
	return parseToken.MouseUser, nil
}
