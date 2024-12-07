package contextkeys

type ContextKey string

const (
	ContextKeyUserID  ContextKey = "userID"
	ContextKeyIsAdmin ContextKey = "isAdmin"
)
