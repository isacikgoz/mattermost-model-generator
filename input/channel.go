//go:generate go run ./cmd/generator model.go

package model

type ChannelType string

type Channel struct {
	id       string      `json:"id"`
	people   []string    `json:"people"`
	name     string      `json:"name" model:"patch,apiPatch"`
	createAt int64       `json:"create_at" model:"patch"`
	cType    ChannelType `json:"type" model:"patch"`
}
