package model

type Player struct {
	ID       string `json:"id"`
	Balance  int    `json:"balance"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
