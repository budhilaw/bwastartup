package mysql

import (
	"belajar-bwa/domain"
	"gorm.io/gorm"
)

type mysqlCampaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) *mysqlCampaignRepository {
	return &mysqlCampaignRepository{db}
}

func (r *mysqlCampaignRepository) FindAll() ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	err := r.db.Find(&campaigns).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *mysqlCampaignRepository) FindByUserID(userID int) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *mysqlCampaignRepository) FindByID(ID int) (domain.Campaign, error) {
	var campaign domain.Campaign
	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", ID).Find(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *mysqlCampaignRepository) Save(campaign domain.Campaign) (domain.Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *mysqlCampaignRepository) Update(campaign domain.Campaign) (domain.Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *mysqlCampaignRepository) CreateImage(campaignImage domain.CampaignImage) (domain.CampaignImage, error) {
	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}
	return campaignImage, nil
}

func (r mysqlCampaignRepository) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	err := r.db.Model(&domain.CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
