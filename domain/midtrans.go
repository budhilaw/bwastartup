package domain

import (
	"github.com/midtrans/midtrans-go"
)

type MidtransTransaction struct {
	ID     int
	Amount int
}

type MidtransUsecase interface {
	GetPaymentURL(transaction MidtransTransaction, user User) (string, *midtrans.Error)
}
