package handler

const (
	ACTION_DEPOSIT  = "DEPOSIT"
	ACTION_WITHDRAW = "WITHDRAW"
)

type History struct {
	Id     string
	Name   string
	Action string
	Value  string
	Date   string
}
