.PHONY:

generate:
	./.github/mockgen-version.sh
	go generate -v ./...

lint: generate
	./.github/lint-version.sh
	golangci-lint run -v ./...

test: generate
	go test -cover ./...

cover: generate
	CGO_ENABLED=1 go test -failfast -count=2 --race -coverprofile=cover.out -coverpkg=./... ./...
	cat cover.out | grep -v "_mock.go" > cover.nomocks.out
	go tool cover -func cover.nomocks.out
