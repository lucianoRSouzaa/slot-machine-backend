package model

type Role string

const (
	PlayerRole Role = "player"
	AdminRole  Role = "admin"
)

type Player struct {
	ID       string `json:"id"`
	Balance  int    `json:"balance"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     Role   `json:"-"`
}
