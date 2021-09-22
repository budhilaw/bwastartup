package usecase

import (
	"belajar-bwa/domain"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserCase(ur domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: ur}
}

func (s *userUsecase) RegisterUser(input domain.RegisterUserInput) (domain.User, error) {
	user := domain.User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passHash)
	user.Role = "user"

	newUser, err := s.userRepo.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *userUsecase) Login(input domain.LoginInput) (domain.User, error) {
	email := input.Email
	password := input.Password

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("No user found on that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userUsecase) IsEmailAvailable(input domain.CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *userUsecase) SaveAvatar(ID int, fileLocation string) (domain.User, error) {
	user, err := s.userRepo.FindByID(ID)
	if err != nil {
		return user, err
	}

	user.AvatarFilename = fileLocation

	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *userUsecase) GetUserByID(ID int) (domain.User, error) {
	user, err := s.userRepo.FindByID(ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("No user found on that ID")
	}

	return user, nil
}
