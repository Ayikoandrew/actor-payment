package storage

import "time"

type CalculatePayment struct {
	PaymentID int64
	Amount    float64
	Currency  string
	Timestamp time.Time
}

type ProcessedPayment struct {
	PaymentID        int64
	InitialPayment   float64
	ProcessedPayment float64
	Status           string
	Timestamp        time.Time
}
