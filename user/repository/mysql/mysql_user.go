package mysql

import (
	"belajar-bwa/domain"
	"gorm.io/gorm"
)

type mysqlUserRepository struct {
	Conn *gorm.DB
}

func NewMysqlUserRepository(Conn *gorm.DB) domain.UserRepository {
	return &mysqlUserRepository{Conn}
}

func (m *mysqlUserRepository) Save(user domain.User) (domain.User, error) {
	err := m.Conn.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *mysqlUserRepository) FindByEmail(email string) (domain.User, error) {
	var user domain.User
	err := m.Conn.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *mysqlUserRepository) FindByID(id int) (domain.User, error) {
	var user domain.User
	err := m.Conn.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *mysqlUserRepository) Update(user domain.User) (domain.User, error) {
	err := m.Conn.Save(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
