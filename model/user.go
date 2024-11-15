package model

type User struct {
	ID       int64  `json:"entityID"`
	Username string `json:"username"`
	Password string `json:"password"`
}
