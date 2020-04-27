//go:generate go run generator.go

package model

type Channel struct {
	id       string `json:"id"`
	name     string `json:"name"`
	createAt int64  `json:"create_at"`
}

type User struct {
	id       string `json:"id"`
	username string `json:"username"`
}
