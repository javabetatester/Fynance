package dashboard

import (
    "Fynance/internal/domain/goal"
    "Fynance/internal/domain/transaction"
)

type Dashboard struct {
    UserId          string
    Name            string
    TotalBalance    float64
    MonthReceipt    float64
    MonthExpense    float64
    FixedExpenses   float64
    TotalInvestment float64
    Goals           []goal.Goal
    Transactions    []transaction.Transaction
}
