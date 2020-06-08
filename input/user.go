package model

type UserType string

type User struct {
	id       string `json:"id"`
	username string `json:"username"`
}
