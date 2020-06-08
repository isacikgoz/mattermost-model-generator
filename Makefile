.PHONY: generate

generate:
	go run ./cmd/generator ./input/channel.go ./input/user.go
	gofmt -s -w ./output/.

migrate_channel:
	go run ./cmd/migrator model Channel  ../mattermost-server