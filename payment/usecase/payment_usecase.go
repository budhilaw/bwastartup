package usecase

import (
	"belajar-bwa/domain"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/spf13/viper"
	"strconv"
)

type midtransUsecase struct {
}

func NewMidtransUsecase() *midtransUsecase {
	return &midtransUsecase{}
}

func (s *midtransUsecase) GetPaymentURL(transaction domain.MidtransTransaction, user domain.User) (string, *midtrans.Error) {
	serverKey := viper.GetString("vendor.midtrans.serverkey")
	midEnv := midtrans.Sandbox

	midsnap := snap.Client{}
	midsnap.New(serverKey, midEnv)

	req := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapRes, err := midsnap.CreateTransaction(req)
	if err != nil {
		return "", err
	}
	return snapRes.RedirectURL, nil
}
