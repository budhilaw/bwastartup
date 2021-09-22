package usecase

import (
	"belajar-bwa/domain"
	"errors"
	"fmt"
	"github.com/gosimple/slug"
)

type campaignUsecase struct {
	campaignRepo domain.CampaignRepository
}

func NewCampaignUsecase(cr domain.CampaignRepository) *campaignUsecase {
	return &campaignUsecase{campaignRepo: cr}
}

func (s *campaignUsecase) GetCampaigns(userID int) ([]domain.Campaign, error) {
	if userID != 0 {
		campaigns, err := s.campaignRepo.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	campaigns, err := s.campaignRepo.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s *campaignUsecase) GetCampaignByID(input domain.GetCampaignDetailInput) (domain.Campaign, error) {
	campaign, err := s.campaignRepo.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *campaignUsecase) CreateCampaign(input domain.CreateCampaignInput) (domain.Campaign, error) {
	campaign := domain.Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate)

	newCampaign, err := s.campaignRepo.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}

func (s *campaignUsecase) UpdateCampaign(inputID domain.GetCampaignDetailInput, inputData domain.CreateCampaignInput) (domain.Campaign, error) {
	campaign, err := s.campaignRepo.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("Not an owner of the campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	updatedCampaign, err := s.campaignRepo.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}
	return updatedCampaign, nil
}

func (s *campaignUsecase) SaveCampaignImage(input domain.CreateCampaignImageInput, fileLocation string) (domain.CampaignImage, error) {
	campaign, err := s.campaignRepo.FindByID(input.CampaignID)
	if err != nil {
		return domain.CampaignImage{}, err
	}

	if campaign.ID == 0 {
		return domain.CampaignImage{}, errors.New("Campaign not found")
	}

	if campaign.UserID != input.User.ID {
		return domain.CampaignImage{}, errors.New("Not an owner of the campaign")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1
		_, err := s.campaignRepo.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return domain.CampaignImage{}, err
		}
	}

	campaignImage := domain.CampaignImage{}
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation

	newCampaignImage, err := s.campaignRepo.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}
