package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IAuthUserUsecase interface {
	SignUp(authUser model.AuthUser) (model.AuthUserResponse, error)
	Login(authUser model.AuthUser) (string, error)
}

type authUserUsecase struct {
	aur repository.IAuthUserRepository
	av  validator.IAuthUserValidator
}

func NewAuthUserUsecase(aur repository.IAuthUserRepository, av validator.IAuthUserValidator) IAuthUserUsecase {
	return &authUserUsecase{aur, av}
}

func (auu *authUserUsecase) SignUp(authUser model.AuthUser) (model.AuthUserResponse, error) {
	if err := auu.av.AuthUserValidate(authUser); err != nil {
		return model.AuthUserResponse{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(authUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.AuthUserResponse{}, err
	}
	newAuthUser := model.AuthUser{
		Email:      authUser.Email,
		Password:   string(hash),
		Department: authUser.Department,
		Name:       authUser.Name,
	}
	if err := auu.aur.CreateAuthUser(&newAuthUser); err != nil {
		return model.AuthUserResponse{}, err
	}
	resAuthUser := model.AuthUserResponse{
		ID:         newAuthUser.ID,
		Email:      newAuthUser.Email,
		Department: newAuthUser.Department,
		Name:       newAuthUser.Name,
	}
	return resAuthUser, nil
}

func (auu *authUserUsecase) Login(authUser model.AuthUser) (string, error) {
	if err := auu.av.AuthUserValidate(authUser); err != nil {
		return "", err
	}
	storedAuthUser := model.AuthUser{}
	if err := auu.aur.GetAuthUserByEmail(&storedAuthUser, authUser.Email); err != nil {
		return "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedAuthUser.Password), []byte(authUser.Password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedAuthUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
