package payment

import (
	"belajar-bwa/user"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/spf13/viper"
	"strconv"
)

type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, *midtrans.Error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, *midtrans.Error) {
	serverKey := viper.GetString("vendor.midtrans.serverkey")
	//clientKey := viper.GetString("vendor.midtrans.clientkey")
	//envMode := viper.GetBool("debug")
	midEnv := midtrans.Sandbox

	//if envMode {
	//	midEnv = midtrans.Sandbox
	//} else {
	//	midEnv = midtrans.Production
	//}

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
