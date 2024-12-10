package models

type GoalProgress struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	TargetAmount float64 `json:"target_amount"`
	SavedAmount  float64 `json:"saved_amount"`
	Progress     float64 `json:"progress"`
}
