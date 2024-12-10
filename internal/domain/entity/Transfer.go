package entity

type Transfer struct {
	ID     int64 `db:"id" json:"id"`
	From   int32 `db:"from" json:"from"`
	To     int32 `db:"to" json:"to"`
	Amount int64 `db:"amount" json:"amount"`
}
