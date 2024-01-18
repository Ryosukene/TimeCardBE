package repository

import (
	"errors"
	"fmt"
	"go-rest-api/model"

	"gorm.io/gorm"
)

type IAuthUserRepository interface {
	GetAuthUserByEmail(authUser *model.AuthUser, email string) error
	CreateAuthUser(authUser *model.AuthUser) error
}

type authUserRepository struct {
	db *gorm.DB
}

func NewAuthUserRepository(db *gorm.DB) IAuthUserRepository {
	return &authUserRepository{db}
}

func (aur *authUserRepository) GetAuthUserByEmail(authUser *model.AuthUser, email string) error {
	if err := aur.db.Where("email = ?", email).First(authUser).Error; err != nil {
		return err
	}
	return nil
}

func (aur *authUserRepository) CreateAuthUser(authUser *model.AuthUser) error {
	// 既に同じメールアドレスを持つユーザーが存在するか確認
	var count int64
	fmt.Println(authUser.Email)
	aur.db.Model(&model.AuthUser{}).Where("email = ?", authUser.Email).Count(&count)

	if count > 0 {
		return errors.New("このメールアドレスは既に使用されています")
	}

	if err := aur.db.Create(authUser).Error; err != nil {
		return err
	}
	return nil
}
