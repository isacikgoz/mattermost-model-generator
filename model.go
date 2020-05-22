//go:generate go run ./cmd/generator model.go

package model

type Channel struct {
	id       string `json:"id"`
	name     string `json:"name" model:"patch,apiPatch"`
	createAt int64  `json:"create_at" model:"patch"`
}

type User struct {
	id       string `json:"id"`
	username string `json:"username"`
}
