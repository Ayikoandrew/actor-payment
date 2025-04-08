package msg

type PublicPayment struct {
	AccountID int64
	Amount    float32
}

type PrivatePayment struct {
	AccountID int64
	Amount    float32
}

type Subscription struct {
	AccountID int64
	Amount    float32
	Period    int32
}

type Error struct {
	AccountID int64
	Reason    string
}
