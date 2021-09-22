package mysql

import (
	"belajar-bwa/domain"
	"gorm.io/gorm"
)

type mysqlTransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *mysqlTransactionRepository {
	return &mysqlTransactionRepository{db}
}

func (r *mysqlTransactionRepository) GetByCampaignID(campaignID int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (r *mysqlTransactionRepository) GetByUserID(userID int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *mysqlTransactionRepository) GetByID(ID int) (domain.Transaction, error) {
	var transaction domain.Transaction

	err := r.db.Where("id = ?", ID).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *mysqlTransactionRepository) Save(transaction domain.Transaction) (domain.Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (r *mysqlTransactionRepository) Update(transaction domain.Transaction) (domain.Transaction, error) {
	err := r.db.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
