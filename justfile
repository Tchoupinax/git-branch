default: build

watch:
  go build -o git-branch *.go || exit 1 && ./git-branch

build:
  go build -o git-branch *.go

test:
  go test ./...

lint:
  go fmt ./... && golangci-lint run
