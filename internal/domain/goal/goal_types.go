package goal

type GoalStatus string

const (
	Active    GoalStatus = "ACTIVE"
	Completed GoalStatus = "COMPLETED"
	Cancelled GoalStatus = "CANCELLED"
)
