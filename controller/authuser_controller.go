package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type IAuthUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error
}

type authUserController struct {
	auu usecase.IAuthUserUsecase
}

func NewAuthUserController(auu usecase.IAuthUserUsecase) IAuthUserController {
	return &authUserController{auu}
}

func (auc *authUserController) SignUp(c echo.Context) error {
	authUser := model.AuthUser{}
	if err := c.Bind(&authUser); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	authUserRes, err := auc.auu.SignUp(authUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, authUserRes)
}

func (auc *authUserController) LogIn(c echo.Context) error {
	authUser := model.AuthUser{}
	if err := c.Bind(&authUser); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	tokenString, err := auc.auu.Login(authUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (auc *authUserController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (auc *authUserController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}
