package usecase

import (
	"belajar-bwa/domain"
	"errors"
	"strconv"
)

type transactionUsecase struct {
	transactionRepo domain.TransactionRepository
	campaignRepo    domain.CampaignRepository
	midtransUsecase domain.MidtransUsecase
}

func NewTransactionUsecase(tr domain.TransactionRepository, cr domain.CampaignRepository, mu domain.MidtransUsecase) *transactionUsecase {
	return &transactionUsecase{
		transactionRepo: tr,
		campaignRepo:    cr,
		midtransUsecase: mu,
	}
}

func (s *transactionUsecase) GetTransactionsByCampaignID(input domain.GetCampaignTransactionsInput) ([]domain.Transaction, error) {
	campaign, err := s.campaignRepo.FindByID(input.ID)
	if err != nil {
		return []domain.Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []domain.Transaction{}, errors.New("Not an owner of the campaign")
	}

	transactions, err := s.transactionRepo.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (s *transactionUsecase) GetTransactionsByUserID(userID int) ([]domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}
	return transactions, err
}

func (s *transactionUsecase) CreateTransaction(input domain.CreateTransactionInput) (domain.Transaction, error) {
	transaction := domain.Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"

	newTransaction, err := s.transactionRepo.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := domain.MidtransTransaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, midErr := s.midtransUsecase.GetPaymentURL(paymentTransaction, input.User)
	if midErr != nil {
		return newTransaction, errors.New(midErr.Message)
	}

	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.transactionRepo.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, err
}

func (s *transactionUsecase) ProcessPayment(input domain.TransactionNotificationInput) error {
	transactionId, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.transactionRepo.GetByID(transactionId)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.transactionRepo.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepo.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err = s.campaignRepo.Update(campaign)
		if err != nil {
			return err
		}
	}
	return nil
}
