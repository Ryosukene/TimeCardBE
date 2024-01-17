package validator

import (
	"go-rest-api/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type IAuthUserValidator interface {
	AuthUserValidate(authUser model.AuthUser) error
}

type authUserValidator struct{}

func NewAuthUserValidator() IAuthUserValidator {
	return &authUserValidator{}
}

func (auv *authUserValidator) AuthUserValidate(authUser model.AuthUser) error {
	return validation.ValidateStruct(&authUser,
		validation.Field(
			&authUser.Email,
			validation.Required.Error("email is required"),
			validation.RuneLength(1, 100).Error("email must be between 1 and 100 characters long"),
			is.Email.Error("invalid email format"),
		),
		validation.Field(
			&authUser.Password,
			validation.Required.Error("password is required"),
			validation.RuneLength(6, 50).Error("password must be between 6 and 50 characters long"),
		),
	)
}
