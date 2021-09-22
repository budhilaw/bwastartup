package domain

import (
	"time"
)

// Entity
type User struct {
	ID             int
	Name           string
	Occupation     string
	Email          string
	PasswordHash   string
	AvatarFilename string
	Role           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Usecase
type UserUsecase interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
}

// Repository
type UserRepository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindByID(id int) (User, error)
	Update(user User) (User, error)
}
