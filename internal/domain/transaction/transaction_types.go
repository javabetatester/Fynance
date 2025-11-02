package transaction

type Types string

const (
	Receipt    Types = "RECEIPT"
	Expense    Types = "EXPENSE"
	Transfer   Types = "TRANSFER"
	Goals      Types = "GOALS"
	Investment Types = "INVESTMENT"
	Withdraw Types = "WITHDRAW"
)
