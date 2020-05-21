.PHONY: generate

generate:
	go generate
	gofmt -s -w ./output/.

migrate_channel:
	go run ./cmd/migrator model Channel  ../mattermost-server/**/*.go