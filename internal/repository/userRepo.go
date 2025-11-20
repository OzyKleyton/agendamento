package repository

import (
	"github.com/OzyKleyton/agendamento-api/internal/model/user"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *user.User) (*user.User, error)
	FindAll() ([]user.User, error)
	FindByID(id uint) (user *user.User, err error)
	FindByEmail(email string) (*user.User, error)
	Update(user *user.User) (*user.User, error)
	Delete(id uint) (user *user.User, err error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) Create(user *user.User) (*user.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) FindAll() (users []user.User, err error) {
	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserRepo) FindByID(id uint) (user *user.User, err error) {
	if err := u.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) FindByEmail(email string) (user *user.User, err error) {
	if err := u.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) Update(user *user.User) (*user.User, error) {
	if err := u.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) Delete(id uint) (user *user.User, err error) {
	if err := u.db.Where("id = ?", id).Delete(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}